package monitoring

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"gitlab.fbs-d.com/dev/go/legacy/helpers/api"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
)

type (
	// IMonitoring интерфейс мониторинга
	IMonitoring interface {
		Counter(*Metric, float64) error
		Inc(*Metric) error
		ExecutionTime(*Metric, func() error) (executionTime float64, err error)
		ObserveExecutionTime(metric *Metric, executionTime time.Duration) (err error)
		Gauge(*Metric, float64) error
		IncGauge(*Metric) error
		DecGauge(*Metric) error
	}

	// Metric метрика для мониторинга
	Metric struct {
		Namespace   string
		Subsystem   string
		Name        string
		ConstLabels prometheus.Labels
	}

	// PrometheusMonitoring реализация IMonitoring для работы с prometheus
	PrometheusMonitoring struct {
		collectors map[string]prometheus.Collector

		lock sync.Mutex

		pushURL      string
		username     string
		password     string
		jobName      string
		instanceName string
		serviceName  string
		subProcess   string

		disableLogPushError bool
		disablePushMetrics  bool
	}

	// Config конфигурация мониторинга
	Config struct {
		PushURL             string `yaml:"pushURL" validate:"nonzero"`
		Username            string `yaml:"username"`
		Password            string `yaml:"password"`
		JobName             string `yaml:"jobName" validate:"nonzero"`
		InstanceName        string `yaml:"instanceName" validate:"nonzero"`
		SubProcess          string `yaml:"subProcess"`
		DisableLogPushError bool   `yaml:"disableLogPushError"`
		DisablePushMetrics  bool   `yaml:"disablePush"`
	}

	collectorType int16
)

const (
	counter collectorType = iota
	histogram
	gauge
)

// NewPrometheusMonitoring создание экземпляра мониторинга prometheus
func NewPrometheusMonitoring(config Config) IMonitoring {
	var m = &PrometheusMonitoring{
		collectors:          map[string]prometheus.Collector{},
		pushURL:             config.PushURL,
		username:            config.Username,
		password:            config.Password,
		jobName:             config.JobName,
		instanceName:        config.InstanceName,
		subProcess:          config.SubProcess,
		disableLogPushError: config.DisableLogPushError,
		disablePushMetrics:  config.DisablePushMetrics,
	}

	go m.push()
	return m
}

// push метод отправляет данные в prometheus раз в 5 секунд
// сам себя перезапускает
func (m *PrometheusMonitoring) push() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(api.SysToken, "monitoring pusher crushed:", err, "restarting...")
			time.Sleep(time.Second)
			go m.push()
		}
	}()
	var ticker = time.NewTicker(time.Second * 5)
	for range ticker.C {
		var pushJob = push.
			New(m.pushURL, m.jobName).
			Gatherer(prometheus.DefaultGatherer).
			Grouping("instance", m.instanceName).
			BasicAuth(m.username, m.password)
		if m.subProcess != "" {
			pushJob.Grouping("subProcess", m.subProcess)
		}
		if !m.disablePushMetrics {
			if err := pushJob.Add(); err != nil && !m.disableLogPushError {
				log.Println(api.SysToken, "error push metrics", m.pushURL, err)
			}
		}
	}
}

// Gauge добавление значения к метрике типа "gauge"
func (m *PrometheusMonitoring) Gauge(metric *Metric, val float64) (err error) {
	var collector prometheus.Collector
	if collector, err = m.collector(gauge, metric); err != nil {
		return err
	}

	var (
		g  prometheus.Gauge
		ok bool
	)
	if g, ok = collector.(prometheus.Gauge); !ok {
		return fmt.Errorf("incorrect collector type. Required prometheus.Gauge, got %T", collector)
	}

	g.Set(val)

	return nil
}

// IncGauge увеличение на 1 метрики типа "gauge"
func (m *PrometheusMonitoring) IncGauge(metric *Metric) (err error) {
	var collector prometheus.Collector
	if collector, err = m.collector(gauge, metric); err != nil {
		return err
	}

	var (
		g  prometheus.Gauge
		ok bool
	)
	if g, ok = collector.(prometheus.Gauge); !ok {
		return fmt.Errorf("incorrect collector type. Required prometheus.Gauge, got %T", collector)
	}

	g.Inc()

	return nil
}

// DecGauge уменьшение на 1 метрики типа "gauge"
func (m *PrometheusMonitoring) DecGauge(metric *Metric) (err error) {
	var collector prometheus.Collector
	if collector, err = m.collector(gauge, metric); err != nil {
		return err
	}

	var (
		g  prometheus.Gauge
		ok bool
	)
	if g, ok = collector.(prometheus.Gauge); !ok {
		return fmt.Errorf("incorrect collector type. Required prometheus.Gauge, got %T", collector)
	}

	g.Dec()

	return nil
}

// Counter добавление значения к метрике типа "счетчик"
func (m *PrometheusMonitoring) Counter(metric *Metric, count float64) (err error) {
	var collector prometheus.Collector
	if collector, err = m.collector(counter, metric); err != nil {
		return err
	}

	var (
		counter prometheus.Counter
		ok      bool
	)
	if counter, ok = collector.(prometheus.Counter); !ok {
		return fmt.Errorf("incorrect collector type. Required prometheus.Counter, got %T", collector)
	}

	counter.Add(count)

	return nil
}

// Inc добавление 1 к метрике типа "счетчик"
func (m *PrometheusMonitoring) Inc(metric *Metric) (err error) {
	return m.Counter(metric, 1)
}

// ExecutionTime добавление информации о времени выполнения функции
func (m *PrometheusMonitoring) ExecutionTime(metric *Metric, h func() error) (executionTime float64, err error) {
	var collector prometheus.Collector
	if collector, err = m.collector(histogram, metric); err != nil {
		return 0, h()
	}

	var (
		histogram prometheus.Histogram
		ok        bool
	)
	if histogram, ok = collector.(prometheus.Histogram); !ok {
		return 0, h()
	}

	var start = time.Now()
	err = h()
	executionTime = time.Since(start).Seconds()
	histogram.Observe(executionTime)
	return
}

func (m *PrometheusMonitoring) ObserveExecutionTime(metric *Metric, executionTime time.Duration) (err error) {
	var collector prometheus.Collector
	if collector, err = m.collector(histogram, metric); err != nil {
		return
	}

	var (
		histogram prometheus.Histogram
		ok        bool
	)
	if histogram, ok = collector.(prometheus.Histogram); !ok {
		return
	}

	histogram.Observe(executionTime.Seconds())
	return
}

// Handler обработчик запроса от prometheus
func (m *PrometheusMonitoring) Handler() http.Handler {
	return promhttp.Handler()
}

// collector добавление метрики в коллекцию
func (m *PrometheusMonitoring) collector(t collectorType, metric *Metric) (collector prometheus.Collector, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	var name = metric.String()
	if _, ok := m.collectors[name]; !ok {
		switch t {
		case counter:
			m.collectors[name] = prometheus.NewCounter(prometheus.CounterOpts{
				Namespace:   metric.Namespace,
				Subsystem:   metric.Subsystem,
				Name:        metric.Name,
				ConstLabels: metric.ConstLabels,
				Help:        "counter " + metric.Name,
			})
		case gauge:
			m.collectors[name] = prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace:   metric.Namespace,
				Subsystem:   metric.Subsystem,
				Name:        metric.Name,
				ConstLabels: metric.ConstLabels,
				Help:        "gauge " + metric.Name,
			})
		case histogram:
			m.collectors[name] = prometheus.NewHistogram(prometheus.HistogramOpts{
				Namespace:   metric.Namespace,
				Subsystem:   metric.Subsystem,
				Name:        metric.Name,
				ConstLabels: metric.ConstLabels,
				Help:        "histogram " + metric.Name,
			})
		}

		if err = prometheus.Register(m.collectors[name]); err != nil {
			return nil, err
		}
	}
	return m.collectors[name], nil
}

// String преобразовывает метрику в строку для отправки
func (m Metric) String() string {
	var (
		keys = make([]string, len(m.ConstLabels))
		i    int
	)
	for k := range m.ConstLabels {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	var h = sha1.New()
	h.Write([]byte(m.Namespace + m.Subsystem + m.Name))

	for _, k := range keys {
		h.Write([]byte(k + m.ConstLabels[k]))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
