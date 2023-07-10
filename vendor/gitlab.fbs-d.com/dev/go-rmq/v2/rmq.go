package rmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/streadway/amqp"
	"gitlab.fbs-d.com/dev/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	deadLetterExchange       = "exchange.dead-letter"
	deadQueueSuffix          = ".dead"
	defaultTimeout           = time.Second * 5
	workerNameRpcConsumer    = "rpcConsumer"
	workerNameQueuePublisher = "queuePublisher"
)

var DefaultLogger = zap.NewNop()

type (
	Rmq interface {
		Connect(func() error) error
		Close() error

		IsConnInitialized() bool
		IsClosed() bool
		State() *State

		Rpc(ctx context.Context, queueName string, msg *amqp.Publishing) (*amqp.Delivery, error)
		ConsumeRpc(queueName string, worker RpcWorkerFunc) error

		Publish(ctx context.Context, queueName string, msg *amqp.Publishing) error
		Consume(queueName string, worker WorkerFunc, skipErrors bool) error
		ConsumeTrialDelay(delay int64, trials int64, exchangeName, routingKey, queueName string, worker WorkerFunc) (err error)
		ConsumeBatch(exchangeName, routingKey, queueName, deadQueueName string, worker WorkerBatchFunc, trials int) error

		QueueBind(name string, key string, exchange string, args amqp.Table) error
		QueueDeclare(name string, args amqp.Table) error

		ExchangeDeclare(exchangeName string, exchangeType string) error
		ExchangeBind(from, to, key string) (err error)
		ExchangePublish(ctx context.Context, exchangeName, exchangeType, routingKey string, msg *amqp.Publishing) error
	}

	State struct {
		CurrentChannelsCount     int
		TotalChannelsCount       int
		TotalChannelsErrorsCount int
	}

	WorkerFunc    func(*amqp.Delivery) error
	RpcWorkerFunc func(*amqp.Delivery) *amqp.Publishing

	queueJob struct {
		name         string
		routing      string
		exchangeType string
		errorsChan   chan error
		message      *amqp.Publishing
	}

	consumerChannel struct {
		deliveryMessage chan *amqp.Delivery
		returnMessage   chan *amqp.Return
	}

	rmq struct {
		conn       *connectionWrapper
		config     *Config
		urlGetter  func() (string, error)
		ctx        context.Context
		cancelFunc func()

		queueJobs   chan queueJob
		returnsChan chan amqp.Return
		respawnChan chan string

		consumerChannels *sync.Map

		rpcResponseQueue string
	}
)

var (
	ErrRequestFailed     = errors.New("request failed")
	ErrRequestTimeout    = errors.New("request timeout")
	ErrPublishTimeout    = errors.New("publish timeout")
	ErrEmptyQueueName    = errors.New("queue name is empty")
	ErrEmptyExchangeName = errors.New("exchange name is empty")
	ErrConnectionClosed  = errors.New("connection closed")
)

func NewRmq(conf *Config) Rmq {
	return NewRmqUrlGetter(
		func() (s string, e error) {
			return conf.Url(), nil
		},
		conf,
	)
}

func NewRmqUrlGetter(getter func() (string, error), conf *Config) Rmq {
	var r = &rmq{
		urlGetter: getter,
		config:    conf,

		rpcResponseQueue: uuid.Must(uuid.NewV4()).String(),

		queueJobs:        make(chan queueJob),
		returnsChan:      make(chan amqp.Return),
		respawnChan:      make(chan string, conf.CountChannelsConsumeRPC+conf.CountChannelsPublisher),
		consumerChannels: &sync.Map{},
	}

	r.ctx, r.cancelFunc = context.WithCancel(context.Background())
	return r
}

func (r *rmq) respawnChannels() {
	for len(r.respawnChan) > 0 {
		<-r.respawnChan
	}

	var connClosed = r.conn.Connection.NotifyClose(make(chan *amqp.Error, 1))
	for {
		select {
		case workerName := <-r.respawnChan:
			select {
			case <-connClosed:
				return
			default:
			}

			var (
				err error
				ch  *amqp.Channel
			)
			if ch, err = r.conn.Channel(); err != nil {
				DefaultLogger.Debug("can't create new amqp-channel", zap.Error(err))
				break
			}
			switch workerName {
			case workerNameRpcConsumer:
				if err = r.rpcConsumer(context.Background(), ch); err != nil {
					DefaultLogger.Debug("can't start new rpcConsumer", zap.Error(err))
					return
				}
			case workerNameQueuePublisher:
				r.queuePublisher(context.Background(), ch)
			default:
				DefaultLogger.Debug("unknown workerName")
			}
		case <-connClosed:
			return
		}
	}
}

// Connect open connection to rmq
func (r *rmq) Connect(handler func() error) (err error) {
	r.conn = newConnectionWrapper()

	// handle starter watchers
	var reconnectHandler = func(conn *connectionWrapper) error {
		DefaultLogger.Debug("start watchers")

		for i := 0; i < r.config.GetConsumeRPCChannelsCount(); i++ {
			var ch *amqp.Channel
			if ch, err = conn.Channel(); err != nil {
				return err
			}

			if err = r.rpcConsumer(context.Background(), ch); err != nil {
				return err
			}
		}

		for i := 0; i < r.config.GetPublishChannelsCount(); i++ {
			var ch *amqp.Channel
			if ch, err = conn.Channel(); err != nil {
				return err
			}

			r.queuePublisher(context.Background(), ch)
		}

		go r.respawnChannels()
		return handler()
	}

	return r.conn.dialUrlGetter(r.urlGetter, reconnectHandler, r.config.TLSClientConfig())
}

// Close stop all consumers, queues and close connection to rmq
func (r *rmq) Close() (err error) {
	r.cancelFunc()
	close(r.queueJobs)
	close(r.respawnChan)
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

func (r *rmq) IsConnInitialized() bool {
	return r.conn != nil
}

func (r *rmq) IsClosed() bool {
	if !r.IsConnInitialized() {
		return true
	}
	return r.conn.IsClosed()
}

func (r *rmq) State() *State {
	return &State{
		CurrentChannelsCount:     r.conn.channelsCount(),
		TotalChannelsCount:       int(r.conn.channelTotalCount),
		TotalChannelsErrorsCount: int(r.conn.channelErrorCount),
	}
}

// queuePublisher new rpc messages watcher and processor
func (r *rmq) queuePublisher(session context.Context, ch *amqp.Channel) {
	go func() {
		defer DefaultLogger.Debug("queuePublisher closed")
		var (
			// Делаем каналы буферизированными потому что в этом цикле выполняем ExchangeDeclare и Publish.
			// При записи в каналы notify/returned, если каналы будут без буфера - ловим дедлок
			notify   = ch.NotifyClose(make(chan *amqp.Error, 1))
			returned = ch.NotifyReturn(make(chan amqp.Return, 1))
		)

		for {
			select {
			case err := <-returned:
				r.returnsChan <- err
			case <-notify:
				r.respawnChan <- workerNameQueuePublisher
				return
			case <-session.Done():
				return
			case j, ok := <-r.queueJobs:
				if !ok {
					select {
					case <-r.ctx.Done():
						return
					default:
					}

					r.respawnChan <- workerNameQueuePublisher
					return
				}
				if len(j.name) > 0 {
					if err := ch.ExchangeDeclare(
						j.name,         // name
						j.exchangeType, // type
						true,           // durable
						false,          // auto-deleted
						false,          // internal
						false,          // no-wait
						nil,            // arguments
					); err != nil {
						j.errorsChan <- errors.Wrap(err, "rmq exchange create error")
						continue
					}
				}

				j.errorsChan <- ch.Publish(j.name, j.routing, true, false, *j.message)
			case <-r.ctx.Done():
				return
			}
		}
	}()
}

// rpcConsumer consume rpc queues
func (r *rmq) rpcConsumer(session context.Context, ch *amqp.Channel) (err error) {
	// create temporary queue
	var responseQueue amqp.Queue
	if responseQueue, err = ch.QueueDeclare(r.rpcResponseQueue, false, true, true, false, nil); err != nil {
		DefaultLogger.Debug(
			"declare queue error",
			zap.String("queue_name", r.rpcResponseQueue),
			zap.Error(err),
		)
		return err
	}

	// listen all received messages
	var msgs <-chan amqp.Delivery
	if msgs, err = ch.Consume(responseQueue.Name, "", true, false, false, false, nil); err != nil {
		DefaultLogger.Debug(
			"error consume response rpc queue",
			zap.String("queue_name", responseQueue.Name),
			zap.Error(err),
		)
		return err
	}

	go func() {
		defer DefaultLogger.Debug("rpcConsumer closed")
		DefaultLogger.Debug("rpc consumer ready", zap.String("queue_name", responseQueue.Name))
		// На всякий случай канал буферизированным, чтобы не словить дедлок, если в цикле будут добавлены
		// операции, которые могут приводить к записи в канал notify
		var notify = ch.NotifyClose(make(chan *amqp.Error, 1))
		for {
			select {
			// process all return messages from channel
			// for example not delivered messages in nonexist queues
			case ret, ok := <-r.returnsChan:
				if !ok {
					return
				}
				// load response chan from temporary "global" rpc watcher
				if responseCh, ok := r.consumerChannels.Load(ret.CorrelationId); ok {
					// return
					responseCh.(*consumerChannel).returnMessage <- &ret
					continue
				}

				DefaultLogger.Debug(
					"no channel",
					zap.String("correlation_id", ret.CorrelationId),
				)
			// receive message from queue
			case d, ok := <-msgs:
				// check empty
				if !ok {
					select {
					case <-r.ctx.Done():
						return
					default:
					}

					r.respawnChan <- workerNameRpcConsumer
					return
				}

				DefaultLogger.Debug("got response", zap.String("correlation_id", d.CorrelationId))

				// load response chan from temporary "global" rpc watcher
				if responseCh, ok := r.consumerChannels.Load(d.CorrelationId); ok {
					// return
					responseCh.(*consumerChannel).deliveryMessage <- &d
					continue
				}

				DefaultLogger.Debug(
					"no channel",
					zap.String("correlation_id", d.CorrelationId),
				)
			case <-notify:
				r.respawnChan <- workerNameRpcConsumer
				return
			case <-session.Done():
				return
			case <-r.ctx.Done():
				return
			}
		}
	}()
	return nil
}

// queuePublish publish messages to watcher queue
func (r *rmq) queuePublish(ctx context.Context, name string, msg *amqp.Publishing) (err error) {
	return r.exchangePublish(ctx, "", "", name, msg)
}

// deadQueuePublish publish delivery messages to dead queue
func (r *rmq) deadQueuePublish(name string, messages []*amqp.Delivery) (err error) {
	for _, d := range messages {
		var ctx, cancelFunc = context.WithTimeout(context.Background(), defaultTimeout)
		if err = r.exchangePublish(ctx, "", "", name, &amqp.Publishing{
			Headers:      d.Headers,
			ContentType:  d.ContentType,
			Body:         d.Body,
			DeliveryMode: amqp.Persistent,
			MessageId:    d.MessageId,
		}); err != nil {
			cancelFunc()
			return err
		}
		cancelFunc()
	}

	return nil
}

// exchangePublish publish messages to watcher queue
func (r *rmq) exchangePublish(ctx context.Context, exchange, exchangeType, routing string, msg *amqp.Publishing) (err error) {
	if r.conn.IsClosed() {
		return ErrConnectionClosed
	}

	DefaultLogger.Debug(
		"publish message",
		zap.String("exchange_name", exchange),
		zap.String("routing", routing),
		zap.String("correlation_id", msg.CorrelationId),
	)

	if _, ok := ctx.Deadline(); !ok {
		var cancelFunc func()
		ctx, cancelFunc = context.WithTimeout(ctx, r.config.GetTimeOut(routing))
		defer cancelFunc()
	}

	var j = queueJob{
		name:         exchange,
		routing:      routing,
		message:      msg,
		exchangeType: exchangeType,
		errorsChan:   make(chan error),
	}

	defer close(j.errorsChan)

	select {
	case r.queueJobs <- j:
		err = <-j.errorsChan
		DefaultLogger.Debug(
			"publish message to queue done",
			zap.String("queue", j.routing),
			zap.String("correlation_id", msg.CorrelationId),
		)
		return err
	case <-ctx.Done():
		return ErrPublishTimeout
	}
}

// ExchangeDeclare create exchange
func (r *rmq) ExchangeDeclare(exchangeName, exchangeType string) (err error) {
	if exchangeName == "" {
		return errors.New("empty exchange name")
	}

	return r.conn.ChannelWrap(func(ch *amqp.Channel) error {
		return ch.ExchangeDeclare(
			exchangeName, // name
			exchangeType, // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // no-wait
			nil,          // arguments
		)
	})
}

// ExchangeBind bind exchange to exchange
func (r *rmq) ExchangeBind(from, to, key string) error {
	if from == "" {
		return errors.New("empty from exchange name")
	}

	if to == "" {
		return errors.New("empty to exchange name")
	}

	if key == "" {
		return errors.New("empty key")
	}

	return r.conn.ChannelWrap(func(ch *amqp.Channel) error {
		return ch.ExchangeBind(to, key, from, false, nil)
	})
}

// Rpc publish message to queue and wait response
func (r *rmq) Rpc(ctx context.Context, queueName string, msg *amqp.Publishing) (d *amqp.Delivery, err error) {
	if queueName == "" {
		return nil, ErrEmptyQueueName
	}

	// set random message CorrelationId for routing from all received messages
	msg.CorrelationId = uuid.Must(uuid.NewV4()).String()
	// set message reply to queue for RPC response
	msg.ReplyTo = r.rpcResponseQueue

	// get message token for debug messages
	var logger = DefaultLogger.With(zap.String("queue", queueName))

	// skip parse if level not debug
	if DefaultLogger.Core().Enabled(zapcore.DebugLevel) {
		logger = logger.With(zap.String("token", getTokenFromMessage(msg)))
	}

	logger.Debug("RPC request [OUT]",
		zap.String("body", string(msg.Body)),
		zap.String("correlation_id", msg.CorrelationId),
		zap.String("response_queue", msg.ReplyTo),
	)

	// create, set and set defer for temporary "global" response watcher.
	// we wait response in ReplyTo queue and return in this deliveryMessage chan by correlation id
	// if rmq can't send message we receive it from return chan with reason
	var consumerChannel = &consumerChannel{
		deliveryMessage: make(chan *amqp.Delivery, 1),
		returnMessage:   make(chan *amqp.Return, 1),
	}
	r.consumerChannels.Store(msg.CorrelationId, consumerChannel)
	defer func() {
		r.consumerChannels.Delete(msg.CorrelationId)
	}()

	if _, ok := ctx.Deadline(); !ok {
		var cancelFunc func()
		ctx, cancelFunc = context.WithTimeout(ctx, r.config.GetTimeOut(queueName))
		defer cancelFunc()
	}

	// publish rpc message
	if err = r.queuePublish(ctx, queueName, msg); err != nil {
		return nil, err
	}

	select {
	// got it
	case d, ok := <-consumerChannel.deliveryMessage:
		if !ok {
			return nil, ErrRequestTimeout
		}
		return d, nil
	case r, ok := <-consumerChannel.returnMessage:
		if !ok {
			return nil, ErrRequestFailed
		}
		return nil, errors.NewWithCode(r.ReplyText, int(r.ReplyCode))
	case <-ctx.Done():
		// timeout
		logger.Debug("RPC request timeout")
		return nil, ErrRequestTimeout
	}
}

// Publish message to queue
func (r *rmq) Publish(ctx context.Context, queueName string, msg *amqp.Publishing) error {
	if queueName == "" {
		return ErrEmptyQueueName
	}

	// extract token only in debug level
	if DefaultLogger.Core().Enabled(zapcore.DebugLevel) {
		DefaultLogger.Debug(
			"message [OUT]",
			zap.String("token", getTokenFromMessage(msg)),
			zap.String("queue", queueName),
			zap.String("body", string(msg.Body)),
		)
	}

	return r.queuePublish(ctx, queueName, msg)
}

func (r *rmq) Consume(queueName string, worker WorkerFunc, skipErrors bool) (err error) {
	if queueName == "" {
		return ErrEmptyQueueName
	}

	var ch *amqp.Channel
	if ch, err = r.conn.Channel(); err != nil {
		return errors.Wrap(err, "error creating channel")
	}
	defer ch.Close()

	if _, err = ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		return err
	}

	if err = ch.Qos(r.config.GetQos(queueName), 0, false); err != nil {
		return errors.Wrap(err, "error setting qos")
	}

	var (
		wg          sync.WaitGroup
		logger      = DefaultLogger.With(zap.String("queue", queueName))
		retryDelays = newConsumerRetryDelays(
			r.config.GetConsumerRetryDelayStart(queueName),
			r.config.GetConsumerRetryDelayStep(queueName),
			r.config.GetConsumerRetryDelayMax(queueName),
		)
	)
	// start consumer goroutines
	for i := 0; i < r.config.GetConsumersCount(queueName); i++ {
		var deliveries <-chan amqp.Delivery
		if deliveries, err = ch.Consume(
			queueName,
			"",    // consumer
			false, // autoAck
			false, // exclusive
			false, // noLocal
			false, // noWait
			nil,   // args
		); err != nil {
			return err
		}

		// start consuming
		wg.Add(1)
		go func(deliveries <-chan amqp.Delivery) {
			var err error
			// consume deliveries
			for delivery := range deliveries {
				delivery := delivery
				logger.Debug(
					"message [IN]",
					zap.String("token", delivery.MessageId),
					zap.String("body", string(delivery.Body)),
				)

				// exec worker func and receive worker error
				if err = worker(&delivery); err != nil {
					logger.Debug(
						"worker error",
						zap.String("token", delivery.MessageId),
						zap.Bool("skipped", skipErrors),
						zap.Error(err),
					)

					// if skip errors just ack message and forget about it
					if skipErrors {
						if err = delivery.Ack(false); err != nil {
							logger.Error(
								"error acking delivery",
								zap.String("token", delivery.MessageId),
								zap.Error(err),
							)
						}
						continue
					}

					// else reject with requeue true after sleeping
					// for monotonically increasing retry delay
					go func(delivery amqp.Delivery) {
						ctx, cancel := context.WithTimeout(r.ctx, retryDelays.getDelay(&delivery))
						defer cancel()
						<-ctx.Done()
						if err := delivery.Reject(true); err != nil {
							logger.Error(
								"error rejecting delivery",
								zap.String("token", delivery.MessageId),
								zap.Error(err),
							)
						}
					}(delivery)
					continue
				}

				// if everything is ok ack message
				retryDelays.maybeClearDelay(&delivery)
				if err = delivery.Ack(false); err != nil {
					logger.Error(
						"error acking delivery",
						zap.String("token", delivery.MessageId),
						zap.Error(err),
					)
				}

				logger.Debug("delivery processed", zap.String("token", delivery.MessageId))
			}
			wg.Done()
		}(deliveries)
	}

	// wait for all consumer goroutines
	// to stop consuming
	wg.Wait()

	return nil
}

// ConsumeTrialDelay listen events and execute with trial and "sleep" by queues
func (r *rmq) ConsumeTrialDelay(
	delay,
	maxTrials int64,
	exchangeName,
	routingKey,
	queueName string,
	worker WorkerFunc,
) (err error) {
	if queueName == "" {
		return ErrEmptyQueueName
	}

	// create channel and set qos
	var ch *amqp.Channel
	if ch, err = r.conn.Channel(); err != nil {
		return errors.Wrap(err, "error create channel")
	}
	defer ch.Close()

	if err = ch.Qos(r.config.GetQos(queueName), 0, false); err != nil {
		return err
	}

	// create dead queue
	var deadQueue = queueName + deadQueueSuffix
	if _, err = ch.QueueDeclare(deadQueue, true, false, false, false, nil); err != nil {
		return err
	}

	var (
		wait   = &sync.WaitGroup{}
		logger = DefaultLogger.With(
			zap.String("exchange", exchangeName),
			zap.String("queue", queueName),
			zap.String("routingKey", routingKey),
		)
	)

	// run consumers in goroutines
	for i := 0; i < r.config.GetConsumersCount(queueName); i++ {
		var msgs <-chan amqp.Delivery
		if msgs, err = ch.Consume(queueName, "", false, false, false, false, nil); err != nil {
			return err
		}

		wait.Add(1)
		go func(work <-chan amqp.Delivery) {
			for d := range work {
				var d = d
				logger.Debug(
					"message [IN]",
					zap.String("token", d.MessageId),
					zap.String("body", string(d.Body)),
				)

				// run worker
				var workerErr = worker(&d)

				// check worker error
				if workerErr != nil {
					logger.Debug(
						"worker error",
						zap.String("token", d.MessageId),
						zap.Error(workerErr),
					)

					// workers trials
					// parse message header and count previous trials
					// if current trial >= max trials
					if parseTrialDelayHeaders(d.Headers, queueName) >= maxTrials-1 {
						// ack message from current queue
						if err = d.Ack(false); err != nil {
							logger.Error(
								"error ack message",
								zap.String("token", d.MessageId),
								zap.Error(err),
							)
							continue
						}

						var ctx, cancelFunc = context.WithTimeout(context.Background(), defaultTimeout)

						// publish message to "dead" queue
						if err = r.queuePublish(ctx, deadQueue, &amqp.Publishing{
							Headers:      d.Headers,
							ContentType:  d.ContentType,
							Body:         d.Body,
							DeliveryMode: amqp.Persistent,
							MessageId:    d.MessageId,
						}); err != nil {
							logger.Error(
								"error send message to dead queue",
								zap.String("dead_queue", deadQueue),
								zap.String("token", d.MessageId),
								zap.Error(err),
							)
						}

						cancelFunc()
						continue
					}

					// if current trial < max trials
					// just reject and message will push to "delay queue" where set ttl and
					// dead letter exchange current queue
					if err = d.Reject(false); err != nil {
						logger.Error(
							"error reject message",
							zap.String("token", d.MessageId),
							zap.Error(err),
						)
					}

					continue
				}

				// if all ok ack message
				if err = d.Ack(false); err != nil {
					logger.Error(
						"error ack message",
						zap.String("token", d.MessageId),
						zap.Error(err),
					)
					continue
				}

				logger.Debug(
					"message processed",
					zap.String("token", d.MessageId),
					zap.Error(err),
				)
			}
			wait.Done()
		}(msgs)
	}
	wait.Wait()
	return nil
}

// ConsumeRpc Consume rpc requests
func (r *rmq) ConsumeRpc(queueName string, worker RpcWorkerFunc) (err error) {
	if queueName == "" {
		return ErrEmptyQueueName
	}

	// create channel and set qos
	var ch *amqp.Channel
	if ch, err = r.conn.Channel(); err != nil {
		return err
	}

	if err = ch.Qos(r.config.GetQos(queueName), 0, false); err != nil {
		return err
	}

	// create temp queue for response
	if _, err = ch.QueueDeclare(queueName, false, true, false, false, map[string]interface{}{
		"x-dead-letter-exchange": deadLetterExchange,
	}); err != nil {
		return err
	}

	var (
		wait   sync.WaitGroup
		logger = DefaultLogger.With(zap.String("queue", queueName))
	)

	for i := 0; i < r.config.GetConsumersCount(queueName); i++ {
		var msgs <-chan amqp.Delivery
		// start consume
		if msgs, err = ch.Consume(queueName, "", false, false, false, false, nil); err != nil {
			continue
		}

		wait.Add(1)
		go func(work <-chan amqp.Delivery) {
			for d := range work {
				var d = d
				logger.Debug("RPC request [IN]",
					zap.String("token", d.MessageId),
					zap.String("correlation_id", d.CorrelationId),
					zap.String("body", string(d.Body)),
				)

				// run rpc worker
				var res = worker(&d)
				// set messages correlation id for correct routing
				res.CorrelationId = d.CorrelationId

				var ctx, cancelFunc = context.WithTimeout(context.Background(), defaultTimeout)
				// publish response
				if err = r.queuePublish(ctx, d.ReplyTo, res); err != nil {
					logger.Error(
						"Error send response",
						zap.String("token", d.MessageId),
						zap.Error(err),
					)
					cancelFunc()
					continue
				}
				cancelFunc()

				// ack current message
				if err = d.Ack(false); err != nil {
					logger.Error(
						"error ack message",
						zap.String("token", d.MessageId),
						zap.Error(err),
					)
					continue
				}

				logger.Debug(
					"message processed",
					zap.String("token", d.MessageId),
					zap.Error(err),
				)
			}
			wait.Done()
		}(msgs)
	}

	wait.Wait()

	return nil
}

// ExchangePublish publish message to exchange
func (r *rmq) ExchangePublish(
	ctx context.Context,
	exchangeName,
	exchangeType,
	routingKey string,
	msg *amqp.Publishing,
) (err error) {
	if exchangeName == "" {
		return ErrEmptyExchangeName
	}

	// skip parse if level not debug
	if DefaultLogger.Core().Enabled(zapcore.DebugLevel) {
		DefaultLogger.Debug(
			"event [OUT]",
			zap.String("token", getTokenFromMessage(msg)),
			zap.String("exchangeName", exchangeName),
			zap.String("routingKey", routingKey),
			zap.String("body", string(msg.Body)),
		)
	}

	// publish message to exchange
	return r.exchangePublish(ctx, exchangeName, exchangeType, routingKey, msg)
}

// Связывает durable очередь с exchange
func (r *rmq) QueueBind(name, key, exchange string, args amqp.Table) (err error) {
	if name == "" {
		return ErrEmptyQueueName
	}

	if exchange == "" {
		return ErrEmptyExchangeName
	}

	return r.conn.ChannelWrap(func(ch *amqp.Channel) error {
		// declare queue
		if _, err = ch.QueueDeclare(name, true, false, false, false, args); err != nil {
			return errors.Wrap(err, "error declare "+name)
		}

		// bind
		return ch.QueueBind(name, key, exchange, false, args)
	})
}

// QueueDeclare создание очереди
func (r *rmq) QueueDeclare(name string, args amqp.Table) (err error) {
	if name == "" {
		return ErrEmptyQueueName
	}

	return r.conn.ChannelWrap(func(ch *amqp.Channel) error {
		if _, err = ch.QueueDeclare(name, true, false, false, false, args); err != nil {
			return errors.Wrap(err, "error declare "+name)
		}

		return nil
	})
}

// QueueDelete delete queue
func (r *rmq) QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (msgs int, err error) {
	if name == "" {
		return 0, ErrEmptyQueueName
	}

	return msgs, r.conn.ChannelWrap(func(ch *amqp.Channel) error {
		msgs, err = ch.QueueDelete(name, ifUnused, ifEmpty, noWait)
		return err
	})
}

func parseTrialDelayHeaders(h amqp.Table, queueName string) (trial int64) {
	var (
		headerQueue, reasonStr string
		data                   amqp.Table
	)
	// check x-death headers
	if add, ok := h["x-death"].([]interface{}); ok {
		// look over x-death reason
		for _, reason := range add {
			// convert to map
			if data, ok = reason.(amqp.Table); !ok {
				continue
			}

			// check queue name from headers
			// fool protection
			if headerQueue, ok = data["queue"].(string); !ok {
				continue
			}

			if headerQueue != fmt.Sprintf("%s.delay", queueName) {
				continue
			}

			// check reason, must be expired
			// expired by TTL from delay queue
			if reasonStr, ok = data["reason"].(string); !ok {
				continue
			}

			if reasonStr != "expired" {
				continue
			}

			// get trial
			if tmpTrial, _ := data["count"].(int64); tmpTrial > 0 {
				trial += tmpTrial
			}
		}
	}
	return trial
}

func getTokenFromMessage(msg *amqp.Publishing) string {
	type TokenizedMessage struct {
		Token string
	}

	var tokentizedMsg = &TokenizedMessage{}
	if err := json.Unmarshal(msg.Body, tokentizedMsg); err != nil {
		return msg.MessageId
	}

	if tokentizedMsg.Token == "" {
		return msg.MessageId
	}

	return tokentizedMsg.Token
}

func newConsumerRetryDelays(
	retryDelayStart,
	retryDelayStep,
	retryDelayMax time.Duration,
) *consumerRetryDelays {
	return &consumerRetryDelays{
		delays:          make(map[string]time.Duration),
		retryDelayStart: retryDelayStart,
		retryDelayStep:  retryDelayStep,
		retryDelayMax:   retryDelayMax,
	}
}

type consumerRetryDelays struct {
	delays          map[string]time.Duration
	mx              sync.RWMutex
	retryDelayStart time.Duration
	retryDelayStep  time.Duration
	retryDelayMax   time.Duration
}

func (cr *consumerRetryDelays) getDelay(delivery *amqp.Delivery) time.Duration {
	cr.mx.RLock()
	prevDelay, isNotFirstRetry := cr.delays[delivery.MessageId]
	cr.mx.RUnlock()

	delay := cr.retryDelayStart
	if isNotFirstRetry {
		delay = prevDelay + cr.retryDelayStep
	}

	if delay > cr.retryDelayMax {
		return cr.retryDelayMax
	}

	cr.mx.Lock()
	cr.delays[delivery.MessageId] = delay
	cr.mx.Unlock()

	return delay
}

func (cr *consumerRetryDelays) maybeClearDelay(delivery *amqp.Delivery) {
	cr.mx.RLock()
	_, clearDelay := cr.delays[delivery.MessageId]
	cr.mx.RUnlock()

	if clearDelay {
		cr.mx.Lock()
		delete(cr.delays, delivery.MessageId)
		cr.mx.Unlock()
	}
}
