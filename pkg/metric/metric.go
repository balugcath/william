package metric

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Interface ...
type Interface interface {
	Register(kind int, name, help string, opts ...string) *Metric
	Add(name string, opts ...interface{}) *Metric
	Set(name string, opts ...interface{}) *Metric
}

const (
	_ = iota
	// GaugeVec ...
	GaugeVec
	// Gauge ...
	Gauge
	// Counter ...
	Counter
	// CounterVec ...
	CounterVec
)

type prom struct {
	kind int
	val  interface{}
}

// Metric ...
type Metric struct {
	port   string
	path   string
	metric map[string]prom
}

// NewMetric ...
func NewMetric(port, path string) *Metric {
	s := &Metric{
		port:   port,
		path:   path,
		metric: make(map[string]prom),
	}
	return s
}

// Register ...
func (s *Metric) Register(kind int, name, help string, opts ...string) *Metric {
	if _, ok := s.metric[name]; ok {
		return s
	}
	switch kind {
	default:
		return s
	case Gauge:
		m := prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: name,
				Help: help,
			},
		)
		prometheus.MustRegister(m)
		s.metric[name] = prom{kind: kind, val: m}

	case GaugeVec:
		m := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: name,
				Help: help,
			},
			opts,
		)
		prometheus.MustRegister(m)
		s.metric[name] = prom{kind: kind, val: m}

	case CounterVec:
		m := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: name,
				Help: help,
			},
			opts,
		)
		prometheus.MustRegister(m)
		s.metric[name] = prom{kind: kind, val: m}

	case Counter:
		m := prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: name,
				Help: help,
			},
		)
		prometheus.MustRegister(m)
		s.metric[name] = prom{kind: kind, val: m}
	}
	return s
}

// Start ...
func (s *Metric) Start() error {
	http.Handle(s.path, promhttp.Handler())
	return http.ListenAndServe(s.port, nil)
}

// Add ...
// Gauge, Counter: opts[0] = float64
// GaugeVec, CounterVec: []opts float64 and strings
func (s *Metric) Add(name string, opts ...interface{}) *Metric {
	// defer func() {
	// 	recover()
	// }()

	v, ok := s.metric[name]
	if !ok {
		return s
	}
	switch v.kind {
	default:
		return s
	case Gauge:
		m, ok := v.val.(prometheus.Gauge)
		if !ok {
			break
		}
		v, ok := opts[0].(float64)
		if !ok {
			break
		}
		m.Add(v)

	case Counter:
		m, ok := v.val.(prometheus.Counter)
		if !ok {
			break
		}
		v, ok := opts[0].(float64)
		if !ok {
			break
		}
		m.Add(v)

	case CounterVec:
		m, ok := v.val.(*prometheus.CounterVec)
		if !ok {
			break
		}
		var (
			opts1 = []string{}
			val   float64
		)

		for i := range opts {
			switch v := opts[i].(type) {
			default:
				continue
			case float64:
				val = v
			case string:
				opts1 = append(opts1, v)
			}
		}
		m.WithLabelValues(opts1...).Add(val)

	case GaugeVec:
		m, ok := v.val.(*prometheus.GaugeVec)
		if !ok {
			break
		}
		var (
			opts1 = []string{}
			val   float64
		)

		for i := range opts {
			switch v := opts[i].(type) {
			default:
				continue
			case float64:
				val = v
			case string:
				opts1 = append(opts1, v)
			}
		}
		m.WithLabelValues(opts1...).Add(val)
	}

	return s
}

// Set ...
// Gauge: opts[0] = float64
// GaugeVec: []opts float64 and strings
func (s *Metric) Set(name string, opts ...interface{}) *Metric {
	// defer func() {
	// 	recover()
	// }()

	v, ok := s.metric[name]
	if !ok {
		return s
	}
	switch v.kind {
	default:
		return s
	case Gauge:
		m, ok := v.val.(prometheus.Gauge)
		if !ok {
			break
		}
		v, ok := opts[0].(float64)
		if !ok {
			break
		}
		m.Set(v)

	case GaugeVec:
		m, ok := v.val.(*prometheus.GaugeVec)
		if !ok {
			break
		}
		var (
			opts1 = []string{}
			val   float64
		)

		for i := range opts {
			switch v := opts[i].(type) {
			default:
				continue
			case float64:
				val = v
			case string:
				opts1 = append(opts1, v)
			}
		}
		m.WithLabelValues(opts1...).Set(val)
	}

	return s
}
