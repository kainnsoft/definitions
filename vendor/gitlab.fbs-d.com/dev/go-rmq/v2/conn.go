// Fork from https://github.com/isayme/go-amqp-reconnect
// refactored
package rmq

import (
	"crypto/tls"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/streadway/amqp"
	"gitlab.fbs-d.com/dev/errors"
	"go.uber.org/zap"
)

var DefaultDelay = delay * time.Second

const delay = 2 // reconnect after delay seconds

// connectionWrapper amqp.connectionWrapper wrapper
type connectionWrapper struct {
	channels []*amqp.Channel
	*amqp.Connection
	lock              *sync.Mutex
	isClosedWrapper   int32
	channelTotalCount uint64
	channelErrorCount uint64
}

func newConnectionWrapper() *connectionWrapper {
	return &connectionWrapper{
		channels: []*amqp.Channel{},
		lock:     &sync.Mutex{},
	}
}

// dialUrlGetter wrap amqp.dial, dial by url function and get a reconnect connection
func (c *connectionWrapper) dialUrlGetter(
	urlGetter func() (string, error),
	handler func(connection *connectionWrapper) error,
	tlsConfig *tls.Config,
) (err error) {
	var url string
	if url, err = urlGetter(); err != nil {
		return err
	}

	c.isClosedWrapper = 0

	var (
		conn             *amqp.Connection
		heartbeatTimeout = time.Second * 20
	)
	if conn, err = amqp.DialConfig(url, amqp.Config{
		Heartbeat:       heartbeatTimeout,
		Locale:          "en_US",
		TLSClientConfig: tlsConfig,
	}); err != nil {
		return err
	}
	c.Connection = conn

	// run handlers
	if err = handler(c); err != nil {
		// got error
		// close all channels
		// close connection
		return c.closeChannels()
	}

	go func() {
		var (
			reason *amqp.Error
			ok     bool
		)
		for {
			// exit this goroutine if closed by developer
			if reason, ok = <-c.Connection.NotifyClose(make(chan *amqp.Error, 1)); ok {
				DefaultLogger.Debug("connection closed", zap.Error(reason))
			}
			DefaultLogger.Debug("after notify event")

			if atomic.LoadInt32(&c.isClosedWrapper) == 1 {
				DefaultLogger.Debug("is closed wrapper")
				return
			}

			// DefaultLogger.Debug("connection closed", zap.Error(reason))
			// again close
			_ = c.closeChannels()

			var trial int
			// try reconnect
		R:
			for {
				DefaultLogger.Debug("try reconnect", zap.Int("trial", trial))
				// get url from getter
				if url, err = urlGetter(); err != nil {
					DefaultLogger.Debug("error get url", zap.Error(err))
					time.Sleep(DefaultDelay)
					continue R
				}

				DefaultLogger.Debug("try dial")

				// try dial
				if conn, err = amqp.DialConfig(url, amqp.Config{
					Heartbeat:       heartbeatTimeout,
					Locale:          "en_US",
					TLSClientConfig: tlsConfig,
				}); err == nil {
					DefaultLogger.Debug("connection success", zap.Error(err))

					trial = 0
					var (
						connOld = (*unsafe.Pointer)(unsafe.Pointer(&c.Connection))
						connNew = unsafe.Pointer(conn)
					)
					// swap pointer
					atomic.SwapPointer(connOld, connNew)
					if err = handler(c); err == nil {
						// connection ok
						// ok run handlers
						DefaultLogger.Debug("reconnect success")
						break R
					}

					// not ok handler, close connection
					_ = conn.Close()
					DefaultLogger.Debug("error create handler", zap.Error(err))
				}

				DefaultLogger.Debug("reconnect failed", zap.Error(err), zap.Int("trial", trial))
				trial++
				// delay for reconnect
				time.Sleep(DefaultDelay)
			}
		}
	}()

	return nil
}

// Channel create channel from amqp connection and store in list
func (c *connectionWrapper) Channel() (channel *amqp.Channel, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if channel, err = c.Connection.Channel(); err != nil {
		atomic.AddUint64(&c.channelErrorCount, 1)
		return nil, err
	}
	_ = channel.Confirm(false)
	atomic.AddUint64(&c.channelTotalCount, 1)
	c.channels = append(c.channels, channel)
	return channel, nil
}

// ChannelWrap create channel from amqp connection, run handler and close it
func (c *connectionWrapper) ChannelWrap(handler func(*amqp.Channel) error) (err error) {
	atomic.AddUint64(&c.channelTotalCount, 1)
	var ch *amqp.Channel
	if ch, err = c.Connection.Channel(); err != nil {
		atomic.AddUint64(&c.channelErrorCount, 1)
		return errors.Wrap(err, "error create channel")
	}
	_ = ch.Confirm(false)
	err = handler(ch)
	_ = ch.Close()
	return err
}

// Close connection and sto reconnect
func (c *connectionWrapper) Close() (err error) {
	if atomic.LoadInt32(&c.isClosedWrapper) == 1 {
		return nil
	}

	atomic.AddInt32(&c.isClosedWrapper, 1)
	return c.closeChannels()
}

// closeChannels close all opened channels
func (c *connectionWrapper) closeChannels() (err error) {
	defer DefaultLogger.Debug("all channel closed")
	c.lock.Lock()
	defer c.lock.Unlock()

	err = c.Connection.Close()
	for _, ch := range c.channels {
		_ = ch.Close()
	}
	c.channels = []*amqp.Channel{}
	return err
}

func (c *connectionWrapper) channelsCount() int {
	c.lock.Lock()
	defer c.lock.Unlock()

	return len(c.channels)
}
