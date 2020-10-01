package metric

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"
)

const (
	path = "/metric"
	port = ":8080"
)

func TestMetric_Add(t *testing.T) {
	type metric struct {
		kind int
		name string
		help string
		opts []string
		data []float64
	}
	type args struct {
		metric []metric
	}
	tests := []struct {
		name string
		want []string
		args args
	}{
		{
			name: "test 1",
			args: args{
				metric: []metric{
					{
						kind: Gauge,
						name: "gauge1",
						data: []float64{10, 11, 12},
					},
					{
						kind: Gauge,
						name: "gauge2",
						data: []float64{20, 21, 22},
					},
					{
						kind: GaugeVec,
						name: "gauge_vec1",
						opts: []string{"one"},
						data: []float64{30, 33},
					},
					{
						kind: Counter,
						name: "counter1",
						opts: []string{"one"},
						data: []float64{30, 33},
					},
					{
						kind: CounterVec,
						name: "counter_vec1",
						opts: []string{"one"},
						data: []float64{30, 33},
					},
				},
			},
			want: []string{
				`gauge1 33`,
				`gauge2 63`,
				`counter1 63`,
				`gauge_vec1{one="one"} 63`,
				`counter_vec1{one="one"} 63`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMetric(port, path)
			for _, v := range tt.args.metric {
				s.Register(v.kind, v.name, v.help, v.opts...)
			}
			for _, v := range tt.args.metric {
				switch v.kind {
				case Gauge, Counter:
					for i := range v.data {
						s.Add(v.name, v.data[i])
					}
				case GaugeVec, CounterVec:
					for i := range v.data {
						s.Add(v.name, []interface{}{v.opts[0], v.data[i]}...)
					}
				}
			}

			go func() {
				s.Start()
			}()

			resp, err := new(http.Client).Get("http://localhost" + port + path)
			if err != nil {
				t.Errorf("TestMetric_Add() error = %v", err)
			}
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("TestMetric_Add() error = %v", err)
			}
			resp.Body.Close()
			t.Log(string(b))
			for _, v := range tt.want {
				matched, err := regexp.Match(v, b)
				if err != nil || !matched {
					t.Errorf("TestMetric_Add() error = %v", err)
				}
			}
		})
	}
}
