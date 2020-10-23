package workpool

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

type r struct{}

func (r) Do(data interface{}) (interface{}, error) {
	return data, nil
}

func TestPool_Request1(t *testing.T) {
	type args struct {
		data          []interface{}
		cntPoolWorker int
		bufferLen     int
	}
	tests := []struct {
		name     string
		args     args
		wantData []interface{}
	}{
		{
			name: "test1",
			args: args{
				cntPoolWorker: 1,
				bufferLen:     0,
				data:          []interface{}{1, 2, 3, 4},
			},
			wantData: []interface{}{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewPool(tt.args.cntPoolWorker, tt.args.bufferLen, r{})
			runtime.Gosched()

			got := make([]interface{}, 0)
			for i := range tt.args.data {
				d, _ := s.Handle(tt.args.data[i])
				runtime.Gosched()
				got = append(got, d)

			}

			if !assert.Equal(t, tt.wantData, got) {
				t.Errorf("Pool.Request1() = %v, want %v", got, tt.wantData)
			}
		})
	}
}

func TestPool_Request2(t *testing.T) {
	type args struct {
		data          []interface{}
		cntPoolWorker int
		bufferLen     int
	}
	tests := []struct {
		name     string
		args     args
		wantData []interface{}
	}{
		{
			name: "test1",
			args: args{
				cntPoolWorker: 1,
				bufferLen:     0,
				data:          []interface{}{1, 2, 3, 4},
			},
			wantData: []interface{}{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewPool(tt.args.cntPoolWorker, tt.args.bufferLen, r{})
			runtime.Gosched()

			got := make([]interface{}, 0)
			for i := range tt.args.data {
				d, _ := s.Handle(tt.args.data[i])
				got = append(got, d)
			}
			if !assert.NotEqual(t, tt.wantData, got) {
				t.Errorf("Pool.Request2() = %v, want %v", got, tt.wantData)
			}
		})
	}
}

func TestPool_Request3(t *testing.T) {
	type args struct {
		cntPoolWorker int
		bufferLen     int
		arrLen        int
	}
	tests := []struct {
		name     string
		args     args
		wantData []interface{}
	}{
		{
			name: "test1",
			args: args{
				cntPoolWorker: 6,
				bufferLen:     0,
				arrLen:        100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewPool(tt.args.cntPoolWorker, tt.args.bufferLen, r{})
			runtime.Gosched()

			for i := 0; i < tt.args.arrLen; i++ {
				d, _ := s.Handle(i)
				runtime.Gosched()
				if d != i {
					t.Errorf("Pool.Request3() = %v, want %v", d, i)
				}
			}
		})
	}
}
