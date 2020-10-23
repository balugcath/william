package queue

import (
	"runtime"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type push struct {
	mtx sync.Mutex
	res []interface{}
}

func (s *push) Do(r interface{}) {
	s.mtx.Lock()
	s.res = append(s.res, r)
	s.mtx.Unlock()
}

type data struct {
	id     int
	series int
}

func (s data) String() string {
	return strconv.Itoa(s.id)
}

func TestQueue_QueueHandle(t *testing.T) {
	type args struct {
		maxWorker int
		push      *push
		req       []data
	}
	tests := []struct {
		name string
		args args
		want []data
	}{
		{
			name: "test1",
			args: args{
				maxWorker: 1,
				push:      &push{res: make([]interface{}, 0)},
				req:       []data{{id: 1}, {id: 2}, {id: 3}},
			},
			want: []data{{id: 1}, {id: 2}, {id: 3}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewQueue(tt.args.maxWorker, tt.args.push)

			for i := range tt.args.req {
				s.Handle(tt.args.req[i])
			}
			runtime.Gosched()

			if !assert.NotEqual(t, tt.args.req, tt.args.push.res) {
				t.Errorf("TestQueue_QueueHandle() = %v, want %v", tt.args.push.res, tt.args.req)
			}
		})
	}
}
