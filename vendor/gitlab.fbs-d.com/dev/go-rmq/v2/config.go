package rmq

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"runtime"
	"strings"
	"time"
)

const (
	defaultConsumerRetryDelayStart = 2 * time.Second
	defaultConsumerRetryDelayStep  = 2 * time.Second
	defaultConsumerRetryDelayMax   = 30 * time.Second

	defaultBatchQos             = 500
	defaultBatchDelay           = time.Second
	defaultChannelWatchersCount = 10

	TLSModeInsecure = "insecure"
)

// Config rmq
type Config struct {
	Host                           string           `yaml:"host"`
	Port                           int              `yaml:"port"`
	Username                       string           `yaml:"username"`
	Password                       string           `yaml:"password"`
	DefaultQos                     int              `yaml:"defaultQos"`
	DefaultBatchQos                int              `yaml:"defaultBatchQos"`
	Qos                            map[string]int   `yaml:"qos"`
	DefaultBatchDelay              int              `yaml:"defaultBatchDelay"`
	BatchDelay                     map[string]int   `yaml:"batchDelay"`
	DefaultBatchConsumersCount     int              `yaml:"defaultBatchConsumersCount"`
	DefaultConsumersCount          int              `yaml:"defaultConsumersCount"`
	ConsumersCount                 map[string]int   `yaml:"consumersCount"`
	DefaultTimeout                 int64            `yaml:"defaultTimeout"`
	Timeouts                       map[string]int64 `yaml:"timeouts"`
	ConsumerRetryDelayStart        map[string]int64 `yaml:"consumerRetryDelayStart"`
	ConsumerRetryDelayStep         map[string]int64 `yaml:"consumerRetryDelayStep"`
	ConsumerRetryDelayMax          map[string]int64 `yaml:"consumerRetryDelayMax"`
	DefaultConsumerRetryDelayStart int64            `yaml:"defaultConsumerRetryDelayStart"`
	DefaultConsumerRetryDelayStep  int64            `yaml:"defaultConsumerRetryDelayStep"`
	DefaultConsumerRetryDelayMax   int64            `yaml:"defaultConsumerRetryDelayMax"`
	CountChannelsPublisher         int              `yaml:"countChannelsPublisher"`
	CountChannelsConsumeRPC        int              `yaml:"countChannelsConsumeRPC"`

	// Use TLSMode=insecure for simple tls connection without domain checking
	TLSMode         string `yaml:"tlsMode"`
	tlsClientConfig *tls.Config
}

func (c *Config) SetTLSClientConfig(cfg *tls.Config) {
	c.TLSMode = ""
	c.tlsClientConfig = cfg
}

func (c *Config) TLSClientConfig() *tls.Config {
	switch c.TLSMode {
	case TLSModeInsecure:
		return &tls.Config{
			// G402: TLS InsecureSkipVerify set true. (gosec)
			InsecureSkipVerify: true, //nolint:gosec
		}
	default:
		return c.tlsClientConfig
	}
}

func (c *Config) Url() string {
	var u = &url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(c.Username, c.Password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
	}
	if c.TLSClientConfig() != nil {
		u.Scheme = "amqps"
	}

	return u.String()
}

func (c *Config) GetConsumersCount(queueName string) int {
	if cnt, ok := c.ConsumersCount[getFormattedQueueName(queueName)]; ok {
		return cnt
	}

	if c.DefaultConsumersCount < 0 {
		return runtime.NumCPU()
	}

	return c.DefaultConsumersCount
}

func (c *Config) GetBatchConsumersCount(queueName string) int {
	if cnt, ok := c.ConsumersCount[getFormattedQueueName(queueName)]; ok {
		return cnt
	}

	if c.DefaultBatchConsumersCount > 0 {
		return c.DefaultBatchConsumersCount
	}

	return runtime.NumCPU()
}

func (c *Config) GetQos(queueName string) int {
	var prefetchCount = c.DefaultQos
	if q, ok := c.Qos[getFormattedQueueName(queueName)]; ok {
		prefetchCount = q
	}
	return prefetchCount
}

func (c *Config) GetBatchQos(queueName string) int {
	if q, ok := c.Qos[getFormattedQueueName(queueName)]; ok {
		return q
	}

	if c.DefaultBatchQos > 0 {
		return c.DefaultBatchQos
	}

	return defaultBatchQos
}

func (c *Config) GetBatchDelay(queueName string) time.Duration {
	if d, ok := c.BatchDelay[getFormattedQueueName(queueName)]; ok {
		return time.Duration(d) * time.Millisecond
	}

	if c.DefaultBatchDelay > 0 {
		return time.Duration(c.DefaultBatchDelay) * time.Millisecond
	}

	return defaultBatchDelay
}

func (c *Config) GetTimeOut(queueName string) time.Duration {
	var timeout = c.DefaultTimeout
	if t, ok := c.Timeouts[getFormattedQueueName(queueName)]; ok {
		timeout = t
	}
	return time.Second * time.Duration(timeout)
}

func (c *Config) GetConsumerRetryDelayStart(queueName string) time.Duration {
	if d, ok := c.ConsumerRetryDelayStart[getFormattedQueueName(queueName)]; ok {
		return time.Duration(d) * time.Millisecond
	}

	if c.DefaultConsumerRetryDelayStart > 0 {
		return time.Duration(c.DefaultConsumerRetryDelayStart) * time.Millisecond
	}

	return defaultConsumerRetryDelayStart
}

func (c *Config) GetConsumerRetryDelayStep(queueName string) time.Duration {
	if d, ok := c.ConsumerRetryDelayStep[getFormattedQueueName(queueName)]; ok {
		return time.Duration(d) * time.Millisecond
	}

	if c.DefaultConsumerRetryDelayStep > 0 {
		return time.Duration(c.DefaultConsumerRetryDelayStep) * time.Millisecond
	}

	return defaultConsumerRetryDelayStep
}

func (c *Config) GetConsumerRetryDelayMax(queueName string) time.Duration {
	if d, ok := c.ConsumerRetryDelayMax[getFormattedQueueName(queueName)]; ok {
		return time.Duration(d) * time.Millisecond
	}

	if c.DefaultConsumerRetryDelayMax > 0 {
		return time.Duration(c.DefaultConsumerRetryDelayMax) * time.Millisecond
	}

	return defaultConsumerRetryDelayMax
}

func (c *Config) GetPublishChannelsCount() int {
	if c.CountChannelsPublisher <= 0 {
		return defaultChannelWatchersCount
	}

	return c.CountChannelsPublisher
}

func (c *Config) GetConsumeRPCChannelsCount() int {
	if c.CountChannelsConsumeRPC <= 0 {
		return defaultChannelWatchersCount
	}

	return c.CountChannelsConsumeRPC
}

// getFormattedQueueName Получаем отформатированное имя очереди, на данный момент только приводим в нижний регистр
// Так как yaml парсер viper приводит ключи массива в нижний регистр, несмотря на то в конфиге они кемелкейс
// Ссылка на обсуждение данной проблемы: https://github.com/spf13/viper/issues/260
// Участок кода из-за которого проблема: https://github.com/spf13/viper/blob/v1.7.1/util.go#L81
func getFormattedQueueName(queueName string) string {
	return strings.ToLower(queueName)
}
