package rtmetric

import (
	"time"

	"github.com/balugcath/william/pkg/metric"
	"github.com/balugcath/william/pkg/types"
)

const (
	reqIntervalSec = 60 * 5
)

type src interface {
	Len() int
}

// RTMetric ...
type RTMetric struct {
}

// NewRTMetric ...
func NewRTMetric(c types.Config, q src, m metric.Interface, t string) *RTMetric {
	s := RTMetric{}
	m.Register(metric.GaugeVec, types.QueueLenMetricName, types.QueueLenMetricHelp, []string{"node", "type"}...)

	go func() {
		for {
			m.Set(types.QueueLenMetricName, []interface{}{c.NodeName, t, float64(q.Len())}...)
			time.Sleep(time.Second * reqIntervalSec)
		}
	}()

	return &s
}
