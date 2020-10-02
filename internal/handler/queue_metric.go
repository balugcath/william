package handler

import (
	"time"

	"github.com/balugcath/william/pkg/metric"
	"github.com/balugcath/william/pkg/types"
)

const (
	reqIntSec = 60 * 5
)

type qm interface {
	Len() int
}

// QueueMetric ...
type QueueMetric struct {
}

// NewQueueMetric ...
func NewQueueMetric(c types.Config, q qm, m metric.Interface, t string) *QueueMetric {
	s := QueueMetric{}
	m.Register(metric.GaugeVec, types.QueueLenMetricName, types.QueueLenMetricHelp, []string{"node", "type"}...)

	go func() {
		for {
			m.Set(types.QueueLenMetricName, []interface{}{c.NodeName, t, float64(q.Len())}...)
			time.Sleep(time.Second * reqIntSec)
		}
	}()

	return &s
}
