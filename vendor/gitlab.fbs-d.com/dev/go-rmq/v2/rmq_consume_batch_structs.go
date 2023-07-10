package rmq

import (
	"errors"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type (
	WorkerBatchFunc func([]*amqp.Delivery) error
	deliveryBatch   struct {
		logger        *zap.Logger
		isDebugLogger bool

		deliveries []*amqp.Delivery
		qos        int
		count      int
		worker     WorkerBatchFunc

		queueName string
		maxTrials int
	}

	batchTimer struct {
		period    time.Duration
		timer     *time.Timer
		isStopped bool
	}
)

// Stop стопим таймер и задаем флаг
func (b *batchTimer) Stop() {
	b.timer.Stop()
	b.isStopped = true
}

// Reset перезапуск счетчика
func (b *batchTimer) Reset() {
	if !b.isStopped {
		// таймер еще работает, ничего не делаем
		return
	}

	// если таймер был остановлен, тогда перезапускаем его
	b.timer.Reset(b.period)
	b.isStopped = false
}

// Add добавление сообщения в пачку
func (d *deliveryBatch) Add(delivery *amqp.Delivery) bool {
	if d.isDebugLogger {
		d.logger.Debug(
			"message added",
			zap.String("token", delivery.MessageId),
			// zap.String("body", string(delivery.Body)),
		)
	}
	d.deliveries = append(d.deliveries, delivery)
	d.count++

	return d.count >= d.qos
}

// Count получение количества сообщений в пачке
func (d *deliveryBatch) Count() int {
	return d.count
}

// reset очистка пачки
func (d *deliveryBatch) reset() {
	d.deliveries = make([]*amqp.Delivery, 0, d.qos)
	d.count = 0
}

// Run запуск обработки пачки
//	- в случае пустой пачки сообщение, выходим
//	- вызываем обработчик пачки
//		- обработчик завершился без ошибок
//			- подтверждаем все сообщения
//			- выходим
//		- при обработке пачки, возникла ошибка
//			- глобальная ошибка (затрагивает все сообщения)
//				- проходимся по всем сообщениям
//					- если у сообщения исчерпаны попытки на повторное выполнение
//						- добавляем сообщения для публикации в dead очередь
//						- подтверждаем его, чтобы выкинуть из текущей очереди
//				- массово отклоняем оставшиеся сообщения
// 			- ошибки для отдельных сообщений (часть выполнилась успешно, а часть нет)
//				- проходимся по ошибочным сообщениям
//					- если у сообщения исчерпаны попытки на повторное выполнение
//						- добавляем сообщения для публикации в dead очередь
//						- подтверждаем его, чтобы выкинуть из текущей очереди
//				- массово принимаем оставшиеся (корректно обработанные) сообщения
//	- очищаем пачку
func (d *deliveryBatch) Run() (deadDeliveries []*amqp.Delivery) {
	if d.count == 0 {
		// защита от пустого списка
		return nil
	}
	defer d.reset() // после выполнения обработки пачки обнуляем расходники

	// вызываем обработчик
	var (
		err       error
		workerErr = d.worker(d.deliveries)
	)
	if workerErr == nil {
		// обработчик завершился без ошибок, подтверждаем все сообщения
		if err = d.deliveries[d.count-1].Ack(true); err != nil {
			d.logger.Error("error ack multiple message", zap.Error(err))
			return nil
		}

		return nil
	}

	deadDeliveries = make([]*amqp.Delivery, 0, d.count)

	var batchErrors *msgErrors
	if !errors.As(workerErr, &batchErrors) {
		d.logger.Debug("common error", zap.Error(workerErr))
		// обработчик отдал общую ошибку для пачки
		var indForNack = -1 // индекс сообщения с которого начинать массовое отклонение
		for ind := 0; ind < len(d.deliveries); ind++ {
			// проверяем что попытки на повторную обработку не исчерпаны
			if int(parseTrialDelayHeaders(d.deliveries[ind].Headers, d.queueName)) < d.maxTrials {
				indForNack = ind // задаем последнее обработанное сообщение, для отклонения
				continue
			}

			// попытки исчерпаны
			// подтверждаем сообщение
			if err = d.deliveries[ind].Ack(false); err != nil {
				d.logger.Error("error ack message", zap.Error(err))
				continue
			}

			// добавляем сообщение в список для публикации в dead очередь
			deadDeliveries = append(deadDeliveries, d.deliveries[ind])
			continue
		}

		if indForNack != -1 {
			// имеются сообщения, которые не попали в dead очередь, отклоняем их
			if err = d.deliveries[indForNack].Nack(true, false); err != nil {
				d.logger.Error("error nack multiple message", zap.Error(err))
				return deadDeliveries
			}
		}

		return deadDeliveries
	}

	var indForAsk = d.count - 1 // задаем последнее сообщение для множественного подтверждения
	// обработчик завершился с ошибками в определенных сообщениях -> отклоняем все ошибочные сообщения
	for _, msgErr := range batchErrors.Errors() {
		d.logger.Debug("batch error", zap.Error(msgErr.Error))
		if indForAsk == msgErr.Index {
			// исключаем принятие уже отклоненного сообщения
			indForAsk--
		}

		if int(parseTrialDelayHeaders(d.deliveries[msgErr.Index].Headers, d.queueName)) < d.maxTrials {
			// попытки на повторную обработку не исчерпаны, отклоняем сообщение
			if err = d.deliveries[msgErr.Index].Reject(false); err != nil {
				d.logger.Error("error reject message", zap.Error(err))
				continue
			}

			continue
		}

		// попытки исчерпаны
		// подтверждаем сообщение
		if err = d.deliveries[msgErr.Index].Ack(false); err != nil {
			d.logger.Error("error ack message", zap.Error(err))
			continue
		}

		// добавляем сообщение в список для публикации в dead очередь
		deadDeliveries = append(deadDeliveries, d.deliveries[msgErr.Index])
		continue
	}

	if indForAsk != -1 {
		// подтверждаем все оставшиеся сообщения
		if err = d.deliveries[indForAsk].Ack(true); err != nil {
			d.logger.Error("error ack multiple message", zap.Error(err))
			return deadDeliveries
		}
	}

	return deadDeliveries
}
