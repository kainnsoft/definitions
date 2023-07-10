package rmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"gitlab.fbs-d.com/dev/go/legacy/exceptions"
	"gitlab.fbs-d.com/dev/go/legacy/helpers"
	"gitlab.fbs-d.com/dev/go/legacy/helpers/api"
	fbs_errors "gitlab.fbs-d.com/dev/go/legacy/helpers/errors"
	"gitlab.fbs-d.com/dev/go/legacy/helpers/monitoring"
	"gitlab.fbs-d.com/dev/go/legacy/i18n"
)

const (
	i18nCategory       = "rmq"
	deadLetterExchange = "exchange.dead-letter"
	defaultMaxChannels = 65535

	MessageInLogPart        = " Message [IN] "
	DelayedMessageInLogPart = " Delayed" + MessageInLogPart
	RpcRequestInLogPart     = " RPC Request [IN] "
)

func init() {
	fbs_errors.OopsT(api.SysToken, i18n.LoadRaw(messages, i18nCategory, ""))
}

type WorkerFunc func(body []byte) excp.IException
type RpcWorkerFunc func(body []byte) (IMessage, excp.IException)

//go:generate mockgen -source rmq.go -destination rmq_mocks.go -package rmq
type IRmq interface {
	Connect() excp.IException
	IsClosed() bool
	Channel() (IChannel, excp.IException)
	ExchangeDeclare(exchangeName string, exchangeType string) excp.IException
	Rpc(queueName string, uuid string, body []byte) ([]byte, excp.IException)
	ConsumeRpc(queueName string, worker RpcWorkerFunc) excp.IException
	Consume(queueName string, worker WorkerFunc) excp.IException
	NonDurableConsume(queueName string, worker WorkerFunc) excp.IException
	Unconsume(queueName string) (ex excp.IException)
	Publish(queueName string, body []byte, headers amqp.Table) excp.IException
	ConsumeDelayed(exchange, queueName string, durable bool, worker WorkerFunc) excp.IException
	PublishDelayed(exchange string, key string, durable bool, body []byte, headers amqp.Table) excp.IException
	ExchangePublish(exchangeName string, exchangeType string, routingKey string, data []byte) excp.IException
	QueueBind(name string, key string, exchange string, noWait bool, args amqp.Table) excp.IException
	NonDurableQueueBind(name string, key string, exchange string, noWait bool, args amqp.Table) excp.IException
	QueueUnbind(name, key, exchange string, args amqp.Table) excp.IException
	QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, excp.IException)
	AddTrimmedFieldForLog(field string)
	TrimFieldsForLog(content []byte) []byte
}

type IChannel interface {
	ExchangeDeclare(name string, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name string, key string, exchange string, noWait bool, args amqp.Table) error
	QueueUnbind(name, key, exchange string, args amqp.Table) error
	QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error)
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	NotifyReturn(ch chan amqp.Return) chan amqp.Return
	Close() error
	Qos(prefetchCount, prefetchSize int, global bool) error
}

type Rmq struct {
	Url                       string
	Conn                      *amqp.Connection
	ErrorFilePath             string
	channelList               sync.Map
	TrimmedFieldsForLog       []string
	TrimmedFieldsForLogRegexp []*regexp.Regexp
	ShortLogQueues            map[string]bool
	ShortLogRegexp            *regexp.Regexp
	MaxChannels               int
	DefaultQos                int
	Qos                       map[string]int
	DefaultConsumersCount     int
	ConsumersCount            map[string]int
	DefaultTimeout            int64
	Timeouts                  map[string]int64
	IsDebugMode               bool
	monitoring                monitoring.IMonitoring
	connectionName            string
	virtualHost               string

	rpcJobs         chan RPCJob
	rpcResponseJobs chan RPCJob
	exchangeJobs    chan ExchangeJob

	rpcResponseQueue string
	consumerChannels sync.Map
	returnsChan      chan amqp.Return
}

type IConsumerChannel interface {
	Send(v amqp.Delivery) bool
	Read() <-chan amqp.Delivery
	Stop()
}

// ConsumerChannel канал для получения ответов на RPC запросы
// на обычном канале невозможно понять, слушает еще кто-то его или нет
// в этой реализации потребитель данных закрывает stopChannel, когда
// он больше не ждет данные, соответственно, если какая-то отправка в канал была
// заблокирована, то закрытие stopChannel отпускает блокировку
type ConsumerChannel struct {
	dataChannel chan amqp.Delivery // канал для данных
	stopChannel chan interface{}   // канал для информирования о закрытии
}

// Send отправка данных в канал
// блокирует поток выполнения либо до успешной отправки, либо до остановки слушателя
// возвращает true, если сообщение было отправлено
// возвращает false, если канал уже закрыт и сообщение доставить не удалось
func (c *ConsumerChannel) Send(v amqp.Delivery) bool {
	select {
	case c.dataChannel <- v:
		return true
	case <-c.stopChannel:
		return false
	}
}

// Read получение данных из канала
func (c *ConsumerChannel) Read() <-chan amqp.Delivery {
	return c.dataChannel
}

// Stop информирование о том, что данные из канала больше не забираются
func (c *ConsumerChannel) Stop() {
	close(c.stopChannel)
}

// NewRmq создает объект Rmq с подключением к Rabbit
func NewRmq(cfg Config, run func(), m monitoring.IMonitoring) (rmq *Rmq, ex excp.IException) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	rmq = new(Rmq)
	rmq.Url = connStr
	rmq.MaxChannels = defaultMaxChannels
	rmq.DefaultConsumersCount = cfg.DefaultConsumersCount
	rmq.ConsumersCount = cfg.ConsumersCount
	rmq.DefaultQos = cfg.DefaultQos
	rmq.Qos = cfg.Qos
	rmq.DefaultTimeout = cfg.DefaultTimeout
	rmq.Timeouts = cfg.Timeouts
	rmq.IsDebugMode = cfg.Debug
	rmq.monitoring = m
	rmq.connectionName = cfg.ConnectionName
	rmq.virtualHost = cfg.VirtualHost

	rmq.rpcJobs = make(chan RPCJob)
	rmq.rpcResponseJobs = make(chan RPCJob)
	rmq.exchangeJobs = make(chan ExchangeJob)
	rmq.rpcResponseQueue, ex = toolkit.Uuid()
	if ex != nil {
		return
	}
	rmq.returnsChan = make(chan amqp.Return)

	for _, trimmed := range cfg.LogTrimmedFields {
		rmq.AddTrimmedFieldForLog(trimmed)
	}

	rmq.ShortLogQueues = make(map[string]bool)
	for _, queueName := range cfg.ShortLogQueues {
		rmq.ShortLogQueues[queueName] = true
	}
	rmq.ShortLogRegexp = regexp.MustCompile(`(?i)"token":"[^"]+"`)

	c := make(chan *amqp.Error)

	// Функция для реконнекта
	go func() {
		for {
			// Ждем ошибку в канал
			rmqStatus := <-c
			if rmqStatus != nil {
				log.Printf("%s RMQ Disconnected. Reason: %s, Code: %d", api.SysToken, rmqStatus.Reason, rmqStatus.Code)
				c = rmq.Reconnect(run)
			}
		}
	}()

	ex = rmq.Connect()
	if ex != nil {
		log.Printf("%s Can't connect to amqp://%s:%s@%s:%d", api.SysToken, cfg.Username, "<hidden>", cfg.Host, cfg.Port)
		c <- new(amqp.Error)
		return
	}
	go rmq.Conn.NotifyClose(c)
	log.Println(api.SysToken, "RMQ Connected")
	return
}

type ExchangeJob struct {
	ErrorsChan   chan excp.IException
	ExchangeName string
	ExchangeType string
	RoutingKey   string
	Body         []byte
	Token        string
}

type RPCJob struct {
	ErrorsChan chan excp.IException
	QueueName  string
	Uuid       string
	Body       []byte
	Token      string
	Timeout    int64
}

func (rmq *Rmq) exchangePublisher(closeChan chan struct{}, exChan chan excp.IException) {
	ch, ex := rmq.Channel()
	if ex != nil {
		exChan <- ex
		return
	}
	log.Println(api.SysToken, "exchangePublisher ready")
	exChan <- nil

	_ = rmq.monitoring.IncGauge(&monitoring.Metric{
		Namespace: "rmq",
		Name:      "exchange_publisher_count",
	})

	defer func() {
		_ = rmq.monitoring.DecGauge(&monitoring.Metric{
			Namespace: "rmq",
			Name:      "exchange_publisher_count",
		})
	}()

L:
	for {
		select {
		case j := <-rmq.exchangeJobs:
			err := ch.ExchangeDeclare(
				j.ExchangeName, // name
				j.ExchangeType, // type
				true,           // durable
				false,          // auto-deleted
				false,          // internal
				false,          // no-wait
				nil,            // arguments
			)
			if err != nil {
				ex = NewRabbitExchangeCreationException(err)
			} else {
				body := j.Body

				err = ch.Publish(
					j.ExchangeName, // exchange
					j.RoutingKey,   // routing key
					false,          // mandatory
					false,          // immediate
					amqp.Publishing{
						ContentType:  "application/json",
						Body:         body,
						DeliveryMode: amqp.Persistent,
						MessageId:    j.Token,
					})
				if err != nil {
					ex = NewRabbitPublishException(err)
				}
			}

			j.ErrorsChan <- ex
		case <-closeChan:
			break L
		}
	}
	log.Println(api.SysToken, "exchangePublisher closed")
}

func (rmq *Rmq) rpcPublisher(closeChan chan struct{}, exChan chan excp.IException) {
	ch, ex := rmq.Channel()
	if ex != nil {
		exChan <- ex
		return
	}
	log.Println(api.SysToken, "rpcPublisher ready")
	exChan <- nil

	rChan := make(chan amqp.Return)
	ch.NotifyReturn(rChan)

	_ = rmq.monitoring.IncGauge(&monitoring.Metric{
		Namespace: "rmq",
		Name:      "rpc_publisher_count",
	})

	defer func() {
		_ = rmq.monitoring.DecGauge(&monitoring.Metric{
			Namespace: "rmq",
			Name:      "rpc_publisher_count",
		})
	}()

L:
	for {
		select {
		case j := <-rmq.rpcJobs:
			ex = nil
			uuid := j.Uuid
			body := j.Body
			queueName := j.QueueName
			token := j.Token

			log.Printf("%s (%s) RPC Request [OUT] %s: %s\n", token, uuid, queueName, string(rmq.TrimFieldsForLog(body)))

			err := ch.Publish("", queueName, true, false, amqp.Publishing{
				ContentType:   "application/json",
				ReplyTo:       rmq.rpcResponseQueue,
				MessageId:     token,
				CorrelationId: uuid,
				Body:          body,
				Expiration:    strconv.FormatInt(j.Timeout*1000, 10),
			})
			if err != nil {
				ex = NewRabbitPublishException(err)
			}

			j.ErrorsChan <- ex
		case r := <-rChan:
			rmq.returnsChan <- r
		case <-closeChan:
			break L
		}
	}
	log.Println(api.SysToken, "rpcPublisher closed")
}

func (rmq *Rmq) rpcResponsePublisher(closeChan chan struct{}, exChan chan excp.IException) {
	ch, ex := rmq.Channel()
	if ex != nil {
		exChan <- ex
		return
	}
	log.Println(api.SysToken, "rpcResponsePublisher ready")
	exChan <- nil

	rChan := make(chan amqp.Return)
	ch.NotifyReturn(rChan)

	_ = rmq.monitoring.IncGauge(&monitoring.Metric{
		Namespace: "rmq",
		Name:      "rpc_response_publisher_count",
	})

	defer func() {
		_ = rmq.monitoring.DecGauge(&monitoring.Metric{
			Namespace: "rmq",
			Name:      "rpc_response_publisher_count",
		})
	}()

L:
	for {
		select {
		case j := <-rmq.rpcResponseJobs:
			uuid := j.Uuid
			body := j.Body
			queueName := j.QueueName
			token := j.Token

			err := ch.Publish("", queueName, false, false, amqp.Publishing{
				ContentType:   "application/json",
				Body:          body,
				CorrelationId: uuid,
				MessageId:     token,
			})
			if err != nil {
				ex = NewRabbitPublishException(err)
			}

			j.ErrorsChan <- ex
		case r := <-rChan:
			rmq.returnsChan <- r
		case <-closeChan:
			break L
		}
	}
	log.Println(api.SysToken, "rpcResponsePublisher closed")
}

func (rmq *Rmq) rpcConsumer(closeChan chan struct{}, exChan chan excp.IException) {
	// канал для получения ответа
	ch, ex := rmq.Channel()
	if ex != nil {
		exChan <- ex
		return
	}

	responseQueue, err := ch.QueueDeclare(rmq.rpcResponseQueue, false, true, true, false, nil)
	if err != nil {
		ex = NewRabbitQueueCreationException(err)
		exChan <- ex
		return
	}

	msgs, err := ch.Consume(responseQueue.Name, "", true, false, false, false, nil)
	if err != nil {
		ex = NewRabbitConsumeException(err)
		exChan <- ex
		return
	}

	log.Println(api.SysToken, "rpcConsumer ready")
	exChan <- nil

	_ = rmq.monitoring.IncGauge(&monitoring.Metric{
		Namespace: "rmq",
		Name:      "rpc_consumer_count",
	})

	defer func() {
		_ = rmq.monitoring.DecGauge(&monitoring.Metric{
			Namespace: "rmq",
			Name:      "rpc_consumer_count",
		})
	}()

L:
	for {
		select {
		case returned := <-rmq.returnsChan:
			log.Printf("%s RPC error %s: %s\n", returned.MessageId, returned.RoutingKey, returned.ReplyText)
			ex := NewRabbitPublishException(errors.New("Message to " + returned.RoutingKey + " not delivered: " + returned.ReplyText))
			if responseCh, ok := rmq.consumerChannels.Load(returned.CorrelationId); ok {
				if !responseCh.(IConsumerChannel).Send(amqp.Delivery{Body: ex.GetJson()}) && rmq.IsDebugMode {
					log.Println(returned.MessageId, "consumer channel is closed")
				}
			}
		case d := <-msgs:
			logToken := d.MessageId
			if logToken == "" {
				logToken = api.SysToken
			}
			if rmq.IsDebugMode {
				log.Println(logToken + " (" + d.CorrelationId + ") RPC Response [IN]")
			}
			if responseCh, ok := rmq.consumerChannels.Load(d.CorrelationId); ok {
				if !responseCh.(IConsumerChannel).Send(d) && rmq.IsDebugMode {
					log.Println(logToken, "consumer channel is closed")
				}
			} else {
				log.Println(logToken, "no channel", d.CorrelationId)
			}
		case <-closeChan:
			break L
		}
	}
	log.Println(api.SysToken, "rpcConsumer closed")
}

func (rmq *Rmq) Reconnect(run func()) chan *amqp.Error {
	rmq.closeConnection()
	runtime.GC()
	// Пытаемся создать новое подключение
	log.Println(api.SysToken, "New RMQ connection")
	i := 1
	for {
		ex := rmq.Connect()
		if ex == nil {
			log.Println(api.SysToken, "Connected")
			break
		}
		log.Println(api.SysToken, "Try reconnect RMQ ", i)
		time.Sleep(5 * time.Second)
		i++
	}
	if run != nil {
		run()
	}
	c := make(chan *amqp.Error)
	go rmq.Conn.NotifyClose(c)
	return c
}

func (rmq *Rmq) Connect() (ex excp.IException) {
	cfg := amqp.Config{
		ChannelMax: rmq.MaxChannels,
		Vhost:      rmq.virtualHost,
		Properties: amqp.Table{
			"connection_name": rmq.connectionName,
		},
	}
	conn, err := amqp.DialConfig(rmq.Url, cfg)
	if err == nil {
		rmq.Conn = conn
	} else {
		ex = NewRabbitConnectionException(err)
	}

	var closeChan chan struct{}
	exChan := make(chan excp.IException, 1)
	for i := 0; i < 5; i++ {
		closeChan = make(chan struct{}, 1)
		go rmq.rpcPublisher(closeChan, exChan)
		ex = <-exChan
		if ex != nil {
			rmq.closeConnection()
			return
		}
		rmq.registerChannel("rpc-publisher-"+strconv.Itoa(i), closeChan)

		closeChan = make(chan struct{}, 1)
		go rmq.rpcResponsePublisher(closeChan, exChan)
		ex = <-exChan
		if ex != nil {
			rmq.closeConnection()
			return
		}
		rmq.registerChannel("rpc-response-publisher-"+strconv.Itoa(i), closeChan)

		closeChan = make(chan struct{}, 1)
		go rmq.exchangePublisher(closeChan, exChan)
		ex = <-exChan
		if ex != nil {
			rmq.closeConnection()
			return
		}
		rmq.registerChannel("exchange-publisher-"+strconv.Itoa(i), closeChan)

		closeChan = make(chan struct{}, 1)
		go rmq.rpcConsumer(closeChan, exChan)

		ex = <-exChan
		if ex != nil {
			rmq.closeConnection()
			return
		}
		rmq.registerChannel("rpc-consumer-"+strconv.Itoa(i), closeChan)
	}

	return
}

// IsClosed returns true if the connection is marked as closed, otherwise false
// is returned.
func (rmq *Rmq) IsClosed() bool {
	return rmq.Conn.IsClosed()
}

func (rmq *Rmq) closeConnection() {
	// Отпускаем все запущенные go-рутины
	log.Println(api.SysToken, "Close RMQ channels")
	rmq.channelList.Range(func(key, value interface{}) bool {
		if ch, ok := value.(chan struct{}); ok {
			ch <- struct{}{}
			return true
		} else {
			return false
		}

	})
	rmq.channelList = sync.Map{}
	rmq.consumerChannels = sync.Map{}
}

func (rmq *Rmq) registerChannel(queueName string, ch chan struct{}) (ex excp.IException) {
	_, ok := rmq.channelList.Load(queueName)
	if ok {
		ex = NewRabbitInternalException(fmt.Errorf("channel for queue %s already exists", queueName))
		return
	}
	rmq.channelList.Store(queueName, ch)
	return
}

// Channel создает канал для обмена сообщениями с попыткой одного реконнекта
func (rmq *Rmq) Channel() (ch IChannel, ex excp.IException) {
	ch, err := rmq.Conn.Channel()
	if err != nil {
		ex = NewRabbitChannelCreationException(err)
	}
	return
}

// Создает exchange
func (rmq *Rmq) ExchangeDeclare(exchangeName string, exchangeType string) (ex excp.IException) {
	if exchangeName == "" {
		ex = NewRabbitExchangeCreationException(errors.New("Empty exchange name"))
		return
	}

	ch, ex := rmq.Channel()
	if ex != nil {
		return
	}
	defer ch.Close()

	err := ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ex = NewRabbitExchangeCreationException(err)
	}

	return
}

// Rpc публикация сообщения в очередь с ожиданием ответа в ответную очередь по уникальному ключу
func (rmq *Rmq) Rpc(queueName string, uuid string, body []byte) (res []byte, ex excp.IException) {
	if queueName == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-queue-name"))
		return
	}
	if uuid == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-uuid"))
		return
	}
	if len(body) == 0 {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-body"))
		return
	}

	if rmq == nil {
		ex = NewRabbitNoConnectionException()
		return
	}

	token := rmq.getTokenFromMessage(body, uuid)
	tokenLog := token
	if rmq.IsDebugMode {
		tokenLog += " (" + uuid + ")"
	}

	var responseBody []byte

	f := func() error {
		errorsChan := make(chan excp.IException)
		defer close(errorsChan)

		responseChan := &ConsumerChannel{make(chan amqp.Delivery), make(chan interface{})}
		defer responseChan.Stop()
		rmq.consumerChannels.Store(uuid, responseChan)
		defer rmq.consumerChannels.Delete(uuid)

		timeout := rmq.DefaultTimeout
		if t, ok := rmq.Timeouts[queueName]; ok {
			timeout = t
		}

		j := RPCJob{
			ErrorsChan: errorsChan,
			QueueName:  queueName,
			Uuid:       uuid,
			Token:      token,
			Body:       body,
			Timeout:    timeout,
		}

		rmq.rpcJobs <- j

		ex = <-errorsChan
		if ex != nil {
			return ex
		}

		// Создаем тикер для определения таймаута
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
		defer cancel()

		func(ctx context.Context) {
			// Ожидаем ответное сообщение
			select {
			case d := <-responseChan.Read():
				responseBody = d.Body
				if e := excp.Parse(d.Body); e != nil {
					ex = e
				} else {
					res = d.Body
				}
			case <-ctx.Done():
				// пришел тик, значит, сообщение еще не обработано
				log.Printf("%s RPC error %s: request timeout\n", token, queueName)
				ex = NewRabbitTimeoutException()
			}
		}(ctx)

		return ex
	}

	executionTime, _ := rmq.pushMetrics(false, queueName, f)
	log.Printf("%s [%.3fs] RPC Response [IN] %s: %s\n", token, executionTime, queueName, string(rmq.TrimFieldsForLog(responseBody)))

	return
}

func (rmq *Rmq) ConsumeRpc(queueName string, worker RpcWorkerFunc) (ex excp.IException) {
	if queueName == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-queue-name"))
		return
	}

	if rmq == nil {
		ex = NewRabbitNoConnectionException()
		return
	}

	consumersQnt := rmq.DefaultConsumersCount
	if cnt, ok := rmq.ConsumersCount[queueName]; ok {
		consumersQnt = cnt
	}

	if consumersQnt == 0 {
		return
	}

	if consumersQnt < 0 {
		consumersQnt = runtime.NumCPU()
	}

	ch, ex := rmq.Channel()
	if ex != nil {
		return
	}
	defer ch.Close()

	// делаем автоудаляющуюся очередь, чтобы отслеживать недоступный сервис
	_, err := ch.QueueDeclare(queueName, false, true, false, false, map[string]interface{}{
		"x-dead-letter-exchange": deadLetterExchange,
	})
	if err != nil {
		ex = NewRabbitQueueCreationException(err)
		return
	}

	forever := make(chan struct{}, 1)
	ex = rmq.registerChannel(queueName, forever)
	if ex != nil {
		return
	}

	prefetchCount := rmq.DefaultQos
	if q, ok := rmq.Qos[queueName]; ok {
		prefetchCount = q
	}

	err = ch.Qos(prefetchCount, 0, false)
	if err != nil {
		ex = NewRabbitChannelCreationException(err)
		return
	}

	for i := 0; i < consumersQnt; i++ {
		msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
		if err != nil {
			ex = NewRabbitConsumeException(err)
			return
		}
		go func(work <-chan amqp.Delivery) {
			for d := range work {
				var responseBody []byte
				token := rmq.getTokenFromMessage(d.Body, d.ReplyTo)

				f := func() (err error) {
					defer func() {
						if ctx := recover(); ctx != nil {
							log.Output(4, fmt.Sprintf("%s RPC Request Error: %v \n Message: %v\n", token, ctx, string(rmq.TrimFieldsForLog(d.Body))))
							if err, ok := ctx.(error); ok {
								rmq.sendRpcResponse(token, d.ReplyTo, NewRabbitInternalException(err).GetJson(), d.CorrelationId)
							}
						}

						e := d.Ack(false)
						if e != nil {
							log.Output(4, fmt.Sprintf("%s RPC Request Ack Error: %v", token, e))
						}
					}()
					log.Println(token + RpcRequestInLogPart + queueName + " : " + rmq.FormatLog(queueName, d.Body))

					res, e := worker(d.Body)
					if e != nil {
						responseBody = e.GetJson()
					} else {
						responseBody = res.GetJson()
					}

					rmq.sendRpcResponse(token, d.ReplyTo, responseBody, d.CorrelationId)

					return e
				}

				executionTime, _ := rmq.pushMetrics(true, queueName, f)
				log.Println(fmt.Sprintf("%s (%s) [%.3fs] RPC Response [OUT] %s: %s\n", token, d.CorrelationId, executionTime, queueName, string(rmq.TrimFieldsForLog(responseBody))))
			}
		}(msgs)
	}
	<-forever

	return
}

func (rmq *Rmq) sendRpcResponse(token, queueName string, body []byte, correlationId string) {
	errorsChan := make(chan excp.IException)
	defer close(errorsChan)

	j := RPCJob{
		ErrorsChan: errorsChan,
		QueueName:  queueName,
		Uuid:       correlationId,
		Token:      token,
		Body:       body,
	}

	rmq.rpcResponseJobs <- j

	ex := <-errorsChan
	if ex != nil {
		_ = log.Output(3, fmt.Sprintf("%s RPC Response send error (publish): %v", token, ex.Error()))
	}
}

func (rmq *Rmq) Consume(queueName string, worker WorkerFunc) (ex excp.IException) {
	if queueName == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-queue-name"))
		return
	}

	if rmq == nil {
		ex = NewRabbitNoConnectionException()
		return
	}

	consumersQnt := rmq.DefaultConsumersCount
	if cnt, ok := rmq.ConsumersCount[queueName]; ok {
		consumersQnt = cnt
	}

	if consumersQnt == 0 {
		return
	}

	if consumersQnt < 0 {
		consumersQnt = runtime.NumCPU()
	}

	ch, ex := rmq.Channel()
	if ex != nil {
		return
	}

	_, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		ex = NewRabbitQueueCreationException(err)
		return
	}

	forever := make(chan struct{}, 1)
	ex = rmq.registerChannel(queueName, forever)
	if ex != nil {
		return
	}

	prefetchCount := rmq.DefaultQos
	if q, ok := rmq.Qos[queueName]; ok {
		prefetchCount = q
	}

	err = ch.Qos(prefetchCount, 0, false)
	if err != nil {
		ex = NewRabbitChannelCreationException(err)
		return
	}

	for i := 0; i < consumersQnt; i++ {
		msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
		if err != nil {
			ex = NewRabbitConsumeException(err)
			return
		}

		go func(work <-chan amqp.Delivery) {
			for d := range work {
				token := rmq.getTokenFromMessage(d.Body, "")

				f := func() (err error) {
					defer func() {
						if ctx := recover(); ctx != nil {
							err = NewRabbitInternalException(fmt.Errorf("panicked: %v", ctx))
						}
					}()

					log.Println(token + MessageInLogPart + queueName + " : " + rmq.FormatLog(queueName, d.Body))

					return worker(d.Body)
				}

				executionTime, err := rmq.pushMetrics(true, queueName, f)

				if err != nil {
					if ex, ok := err.(excp.IException); ok {
						log.Printf("%s [%.3fs] Message [IN] %s Error: %s (%d)", token, executionTime, queueName, ex.Error(), ex.Code())
					} else {
						log.Printf("%s [%.3fs] Message [IN] %s Error: %s", token, executionTime, queueName, err)
					}
					time.Sleep(5 * time.Second)
					_ = d.Reject(true)
				} else {
					log.Printf("%s [%.3fs] Message [IN] %s: Processed", token, executionTime, queueName)
					_ = d.Ack(false)
				}
			}
		}(msgs)
	}
	<-forever
	return
}

// Метод отпускает  consume рутину и удаляет канал из channelList
func (rmq *Rmq) Unconsume(queueName string) (ex excp.IException) {
	if queueName == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-queue-name"))
		return
	}

	if channel, ok := rmq.channelList.Load(queueName); ok {
		if ch, ok := channel.(chan struct{}); ok {
			ch <- struct{}{}
			rmq.channelList.Delete(queueName)
		}
	}
	return
}

func (rmq *Rmq) Publish(queueName string, data []byte, headers amqp.Table) (ex excp.IException) {
	defer func() {
		if r := recover(); r != nil {
			filename := rmq.ErrorFilePath + strconv.Itoa(int(time.Now().Unix()))
			var f *os.File
			var err error
			if _, exists := os.Stat(filename); exists != nil {
				f, err = os.Create(filename)
			} else {
				f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			}
			if err == nil {
				defer f.Close()
			}

			_, err = f.Write([]byte(queueName + "\n"))
			_, err = f.Write(data)
			if err != nil {
				log.Println(api.SysToken, "Cant write to file "+filename+" "+string(rmq.TrimFieldsForLog(data)))
			}
		}
	}()
	if queueName == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-queue-name"))
		return
	}

	if rmq == nil {
		ex = NewRabbitNoConnectionException()
		return
	}

	ch, ex := rmq.Channel()
	if ex != nil {
		return
	}
	defer ch.Close()

	_, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		ex = NewRabbitQueueCreationException(err)
		return
	}

	token := rmq.getTokenFromMessage(data, "")
	log.Output(2, fmt.Sprintf("%s Message [OUT] %s: %s\n", token, queueName, string(rmq.TrimFieldsForLog(data))))

	err = ch.Publish("", queueName, false, false, amqp.Publishing{
		Headers:      headers,
		ContentType:  "application/json",
		Body:         data,
		DeliveryMode: amqp.Persistent,
		MessageId:    token,
	})
	if err != nil {
		ex = NewRabbitPublishException(err)
	}

	return
}

// Consume non-durable очереди (если очереди нет - она объявляется с параметром durable: false, autoDelete: true)
func (rmq *Rmq) NonDurableConsume(queueName string, worker WorkerFunc) (ex excp.IException) {
	if queueName == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-queue-name"))
		return
	}

	consumersQnt := rmq.DefaultConsumersCount
	if cnt, ok := rmq.ConsumersCount[queueName]; ok {
		consumersQnt = cnt
	}

	if consumersQnt == 0 {
		return
	}

	if consumersQnt < 0 {
		consumersQnt = runtime.NumCPU()
	}

	ch, ex := rmq.Channel()
	if ex != nil {
		return
	}

	_, err := ch.QueueDeclare(queueName, false, true, false, false, nil)
	if err != nil {
		ex = NewRabbitQueueCreationException(err)
		return
	}

	ch.Close()

	forever := make(chan struct{}, 1)
	ex = rmq.registerChannel(queueName, forever)
	if ex != nil {
		return
	}

	prefetchCount := rmq.DefaultQos
	if q, ok := rmq.Qos[queueName]; ok {
		prefetchCount = q
	}

	for i := 0; i < consumersQnt; i++ {
		var ch IChannel
		ch, ex = rmq.Channel()
		if ex != nil {
			return
		}
		defer ch.Close()
		err = ch.Qos(prefetchCount, 0, false)
		if err != nil {
			ex = NewRabbitChannelCreationException(err)
			return
		}

		msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
		if err != nil {
			ex = NewRabbitConsumeException(err)
			return
		}

		go func(work <-chan amqp.Delivery) {
			for d := range work {
				func() {
					token := rmq.getTokenFromMessage(d.Body, "")

					defer func() {
						if ctx := recover(); ctx != nil {
							log.Output(2, fmt.Sprintf("%s Message [IN] %s Error: %v \n Message: %v", token, queueName, ctx, string(rmq.TrimFieldsForLog(d.Body))))
							time.Sleep(5 * time.Second)
							d.Reject(true)
						}
					}()

					log.Println(token + MessageInLogPart + queueName + " : " + rmq.FormatLog(queueName, d.Body))

					startTime := time.Now()
					e := worker(d.Body)
					elapsed := toolkit.SinceTimeMilliSeconds(startTime)
					if e != nil {
						log.Output(2, fmt.Sprintf("%s [%.2fms] Message [IN] %s Error: %s (%d)", token, elapsed, queueName, e.Error(), e.Code()))
						time.Sleep(5 * time.Second)
						d.Reject(true)
					} else {
						log.Output(2, fmt.Sprintf("%s [%.2fms] Message [IN] %s: Processed", token, elapsed, queueName))
						d.Ack(false)
					}
				}()
			}
		}(msgs)
	}
	<-forever
	return
}

func (rmq *Rmq) ConsumeDelayed(exchangeName string, queueName string, durable bool, worker WorkerFunc) (ex excp.IException) {
	if exchangeName == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-exchange-name"))
		return
	}

	if rmq == nil {
		ex = NewRabbitNoConnectionException()
		return
	}

	consumersQnt := rmq.DefaultConsumersCount
	if cnt, ok := rmq.ConsumersCount[queueName]; ok {
		consumersQnt = cnt
	}

	if consumersQnt == 0 {
		return
	}

	if consumersQnt < 0 {
		consumersQnt = runtime.NumCPU()
	}

	ch, ex := rmq.Channel()
	if ex != nil {
		return
	}

	err := ch.ExchangeDeclare(exchangeName, "x-delayed-message", durable, !durable, false, false, amqp.Table{"x-delayed-type": "direct"})
	if err != nil {
		ex = NewRabbitExchangeCreationException(err)
		return
	}

	_, err = ch.QueueDeclare(queueName, durable, !durable, false, false, nil)
	if err != nil {
		ex = NewRabbitQueueCreationException(err)
		return
	}

	err = ch.QueueBind(queueName, queueName, exchangeName, false, nil)
	if err != nil {
		ex = NewRabbitQueueBindException(err)
		return
	}

	ch.Close()

	forever := make(chan struct{}, 1)
	ex = rmq.registerChannel(queueName, forever)
	if ex != nil {
		return
	}

	prefetchCount := rmq.DefaultQos
	if q, ok := rmq.Qos[queueName]; ok {
		prefetchCount = q
	}

	for i := 0; i < consumersQnt; i++ {
		var ch IChannel
		ch, ex = rmq.Channel()
		if ex != nil {
			return
		}
		defer ch.Close()
		err = ch.Qos(prefetchCount, 0, false)
		if err != nil {
			ex = NewRabbitChannelCreationException(err)
			return
		}

		msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
		if err != nil {
			ex = NewRabbitConsumeException(err)
			return
		}

		go func(work <-chan amqp.Delivery) {
			for d := range work {
				token := rmq.getTokenFromMessage(d.Body, "")

				f := func() (err error) {
					defer func() {
						if ctx := recover(); ctx != nil {
							err = NewRabbitInternalException(fmt.Errorf("panicked: %v", ctx))
						}
					}()

					log.Println(token + DelayedMessageInLogPart + queueName + " : " + rmq.FormatLog(queueName, d.Body))

					return worker(d.Body)
				}

				executionTime, err := rmq.pushMetrics(true, queueName, f)

				if err != nil {
					if ex, ok := err.(excp.IException); ok {
						log.Printf("%s [%.3fs] Delayed Message [IN] %s/%s Error: %s (%d)", token, executionTime, exchangeName, queueName, ex.Error(), ex.Code())
					} else {
						log.Printf("%s [%.3fs] Delayed Message [IN] %s/%s Error: %s", token, executionTime, exchangeName, queueName, err)
					}
					time.Sleep(5 * time.Second)
					_ = d.Reject(true)
				} else {
					log.Printf("%s [%.3fs] Delayed message [IN] %s/%s: Processed", token, executionTime, exchangeName, queueName)
					_ = d.Ack(false)
				}
			}
		}(msgs)
	}
	<-forever
	return
}

func (rmq *Rmq) PublishDelayed(exchangeName string, key string, durable bool, data []byte, headers amqp.Table) (ex excp.IException) {
	defer func() {
		if r := recover(); r != nil {
			filename := rmq.ErrorFilePath + strconv.Itoa(int(time.Now().Unix()))
			var f *os.File
			var err error
			if _, exists := os.Stat(filename); exists != nil {
				f, err = os.Create(filename)
			} else {
				f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			}
			if err == nil {
				defer f.Close()
			}

			_, err = f.Write([]byte(exchangeName + "\n"))
			_, err = f.Write(data)
			if err != nil {
				log.Println(api.SysToken, "Cant write to file "+filename+" "+string(rmq.TrimFieldsForLog(data)))
			}
		}
	}()

	if exchangeName == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-exchange-name"))
		return
	}

	if rmq == nil {
		ex = NewRabbitNoConnectionException()
		return
	}

	ch, ex := rmq.Channel()
	if ex != nil {
		return
	}
	defer ch.Close()

	err := ch.ExchangeDeclare(exchangeName, "x-delayed-message", durable, !durable, false, false, amqp.Table{"x-delayed-type": "direct"})
	if err != nil {
		ex = NewRabbitExchangeCreationException(err)
		return
	}

	token := rmq.getTokenFromMessage(data, "")
	log.Output(2, fmt.Sprintf("%s Delayed message [OUT] %s/%s: %s\n", token, exchangeName, key, string(rmq.TrimFieldsForLog(data))))

	err = ch.Publish(exchangeName, key, false, false, amqp.Publishing{
		Headers:      headers,
		ContentType:  "application/json",
		Body:         data,
		DeliveryMode: amqp.Persistent,
		MessageId:    token,
	})
	if err != nil {
		ex = NewRabbitPublishException(err)
	}

	return
}

func (rmq *Rmq) ExchangePublish(exchangeName string, exchangeType string, routingKey string, data []byte) (ex excp.IException) {
	if exchangeName == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-exchange-name"))
		return
	}

	token := rmq.getTokenFromMessage(data, "")
	log.Output(2, fmt.Sprintf("%s Event [OUT] %s/%s: %s\n", token, exchangeName, routingKey, string(rmq.TrimFieldsForLog(data))))

	errorsChan := make(chan excp.IException)
	defer close(errorsChan)

	j := ExchangeJob{
		ErrorsChan:   errorsChan,
		ExchangeName: exchangeName,
		ExchangeType: exchangeType,
		RoutingKey:   routingKey,
		Token:        token,
		Body:         data,
	}

	rmq.exchangeJobs <- j

	ex = <-errorsChan
	return
}

// Связывает durable очередь с exchange
func (rmq *Rmq) QueueBind(name string, key string, exchange string, noWait bool, args amqp.Table) (ex excp.IException) {
	if name == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-queue-name"))
		return
	}

	if exchange == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-exchange-name"))
		return
	}

	if rmq == nil {
		ex = NewRabbitNoConnectionException()
		return
	}

	// если количество консьюмеров для очереди указано 0, то не создаем очередь
	consumersQnt := rmq.DefaultConsumersCount
	if cnt, ok := rmq.ConsumersCount[name]; ok {
		consumersQnt = cnt
	}

	if consumersQnt == 0 {
		return
	}

	ch, ex := rmq.Channel()
	if ex != nil {
		return
	}
	defer ch.Close()

	_, err := ch.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		ex = NewRabbitQueueCreationException(err)
		return
	}

	err = ch.QueueBind(name, key, exchange, false, nil)
	if err != nil {
		ex = NewRabbitQueueBindException(err)
	}

	return
}

// Отвязывает очередь  от exchange
func (rmq *Rmq) QueueUnbind(name, key, exchange string, args amqp.Table) (ex excp.IException) {
	if name == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-queue-name"))
		return
	}

	if exchange == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-exchange-name"))
		return
	}

	ch, ex := rmq.Channel()
	if ex != nil {
		return
	}
	defer ch.Close()

	err := ch.QueueUnbind(name, key, exchange, args)
	if err != nil {
		ex = NewRabbitQueueUnbindException(err)
	}

	return
}

// Связывает non-durable очередь с exchange
func (rmq *Rmq) NonDurableQueueBind(name string, key string, exchange string, noWait bool, args amqp.Table) (ex excp.IException) {
	if name == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-queue-name"))
		return
	}

	if exchange == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-exchange-name"))
		return
	}

	consumersQnt := rmq.DefaultConsumersCount
	if cnt, ok := rmq.ConsumersCount[name]; ok {
		consumersQnt = cnt
	}

	if consumersQnt == 0 {
		return
	}

	ch, ex := rmq.Channel()
	if ex != nil {
		return
	}
	defer ch.Close()

	_, err := ch.QueueDeclare(name, false, true, false, false, nil)
	if err != nil {
		ex = NewRabbitQueueCreationException(err)
		return
	}

	err = ch.QueueBind(name, key, exchange, noWait, nil)
	if err != nil {
		ex = NewRabbitQueueBindException(err)
	}

	return
}

// Удаляет очередь
func (rmq *Rmq) QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (msgs int, ex excp.IException) {
	if name == "" {
		ex = NewRabbitInvalidParamException(errors.New("error.empty-queue-name"))
		return
	}

	ch, ex := rmq.Channel()
	if ex != nil {
		return
	}
	defer ch.Close()
	ex = rmq.Unconsume(name)
	if ex != nil {
		return
	}

	msgs, err := ch.QueueDelete(name, ifUnused, ifEmpty, noWait)
	if err != nil {
		ex = NewRabbitQueueDeleteException(err)
	}
	return
}

// Добавляет поле, значение которого должно обрезаться в логах
func (rmq *Rmq) AddTrimmedFieldForLog(field string) {
	rmq.TrimmedFieldsForLog = append(rmq.TrimmedFieldsForLog, field)
	rmq.TrimmedFieldsForLogRegexp = append(rmq.TrimmedFieldsForLogRegexp, regexp.MustCompile(fmt.Sprintf(`\\*"%s\\*":\s*\\*("(?:[^"\\]|\\.)*\\*"|\[\"[^].]*\"\])`, field)))
}

// Обрезает значения полей для лога
func (rmq *Rmq) TrimFieldsForLog(content []byte) []byte {
	for i, re := range rmq.TrimmedFieldsForLogRegexp {
		content = re.ReplaceAll(content, []byte(fmt.Sprintf(`"%s":"TRIMMED_CONTENT"`, rmq.TrimmedFieldsForLog[i])))
	}

	return content
}

// FormatLog  вырезает либо весь лог за исключением токена, либо определённые поля
func (rmq *Rmq) FormatLog(queueName string, content []byte) string {
	if rmq.ShortLogQueues[queueName] {
		cut := string(rmq.ShortLogRegexp.Find(content))
		if len(cut) > 0 {
			return "{" + cut + "}"
		}
	}
	return string(rmq.TrimFieldsForLog(content))
}

func (rmq *Rmq) getTokenFromMessage(body []byte, uuid string) string {
	msg := new(TokenizedMessage)
	err := json.Unmarshal(body, msg)
	if err != nil {
		return uuid
	}

	return msg.Token
}

func (rmq *Rmq) pushMetrics(isServer bool, queueName string, f func() error) (executionTime float64, err error) {
	executionTime, err = rmq.monitoring.ExecutionTime(&monitoring.Metric{
		Namespace: "rmq",
		Name:      "execution_time",
		ConstLabels: map[string]string{
			"queue":     queueName,
			"is_server": fmt.Sprintf("%t", isServer),
		},
	}, f)

	if err != nil {
		_ = rmq.monitoring.Inc(
			&monitoring.Metric{
				Namespace: "rmq",
				Name:      "errors",
				ConstLabels: map[string]string{
					"queue":     queueName,
					"code":      strconv.Itoa(err.(excp.IException).Code()),
					"is_server": fmt.Sprintf("%t", isServer),
				},
			},
		)
	}

	return
}
