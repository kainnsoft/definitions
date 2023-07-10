package rmq

import (
	"sync"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ConsumeBatch слушатель, который обрабатывает сообщения пачками
func (r *rmq) ConsumeBatch(
	exchangeName,
	routingKey,
	queueName,
	deadQueueName string,
	worker WorkerBatchFunc,
	trials int,
) (err error) {
	var l = DefaultLogger.With(
		zap.String("exchange_name", exchangeName),
		zap.String("routing", routingKey),
		zap.String("queueName", queueName),
		zap.String("deadQueueName", deadQueueName),
	)

	if queueName == "" {
		l.Error(ErrEmptyQueueName.Error())
		return ErrEmptyQueueName
	}

	var (
		batchDelay     = r.config.GetBatchDelay(queueName)
		consumersCount = r.config.GetBatchConsumersCount(queueName)
		qos            = r.config.GetBatchQos(queueName)
		wg             = &sync.WaitGroup{}
	)
	l = l.With(zap.Int("qos", qos), zap.Duration("batchDelay", batchDelay))

	// запускаем каналы в goroutines
	wg.Add(consumersCount)
	for i := 0; i < consumersCount; i++ {
		go func(i int) {
			var l = l.With(zap.Int("channel", i))
			// создаем канал
			var ch *amqp.Channel
			if ch, err = r.conn.Channel(); err != nil {
				l.Error("error create channel", zap.Error(err))
				return
			}
			var channelClosed bool
			defer func() {
				if !channelClosed {
					if err = ch.Close(); err != nil {
						l.Error("error close channel", zap.Error(err))
					}
				}
			}()

			// добавляем количество обрабатываемых сообщений (необходимо для обработки пачек сообщений)
			if err = ch.Qos(qos, 0, false); err != nil {
				l.Error("error set Qos", zap.Error(err))
				return
			}

			var msgs <-chan amqp.Delivery
			if msgs, err = ch.Consume(queueName, "", false, false, false, false, nil); err != nil {
				l.Error("error set consume", zap.Error(err))
				return
			}

			var (
				batch = &deliveryBatch{
					logger:        l,
					isDebugLogger: l.Core().Enabled(zapcore.DebugLevel),
					deliveries:    make([]*amqp.Delivery, 0, qos),
					qos:           qos,
					worker:        worker,
					queueName:     queueName,
					maxTrials:     trials,
				}
				batchT = batchTimer{
					period: batchDelay,
					timer:  time.NewTimer(0),
				}
			)
			for {
				select {
				case <-r.ctx.Done():
					channelClosed = true
					batchT.timer.Stop() // останавливаем таймер
					wg.Done()           // завершаем goroutine
					return
				case <-batchT.timer.C:
					// обработка по таймеру
					if deadDeliveries := batch.Run(); deadDeliveries != nil {
						// публикуем сообщения для dead очереди
						_ = r.deadQueuePublish(deadQueueName, deadDeliveries)
					}
					batchT.Stop() // останавливаем таймер
				case i, ok := <-msgs:
					if !ok {
						return
					}
					if isRun := batch.Add(&i); isRun { // пачка заполнилась
						batchT.Stop() // останавливаем таймер
						// вызываем обработку
						if deadDeliveries := batch.Run(); deadDeliveries != nil {
							// публикуем сообщения для dead очереди
							_ = r.deadQueuePublish(deadQueueName, deadDeliveries)
						}

						continue
					}
					// обновляем таймер, если нужно
					batchT.Reset()
				}
			}
		}(i)
	}
	wg.Wait()
	return nil
}
