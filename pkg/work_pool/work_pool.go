package workpool

import (
	"errors"
)

// Interface ...
type Interface interface {
	Do(interface{}) (interface{}, error)
}

type job struct {
	err  error
	data interface{}
}

var (
	// ErrCntWorker ...
	ErrCntWorker = errors.New("count workers pool must be greater than zero")
	// ErrBufferFull ...
	ErrBufferFull = errors.New("buffer full")
)

// Pool ...
type Pool struct {
	ch chan chan job
}

// NewPool ...
func NewPool(cntPoolWorker int, bufferLen int, req Interface) (*Pool, error) {
	if cntPoolWorker == 0 {
		return nil, ErrCntWorker
	}
	s := Pool{
		ch: make(chan chan job, bufferLen),
	}

	for i := 0; i < cntPoolWorker; i++ {
		go func(id int) {
			for {
				reqCh := <-s.ch
				reqJob := <-reqCh
				resp, err := req.Do(reqJob.data)
				if err != nil {
					reqCh <- job{data: resp, err: err}
					continue
				}
				reqCh <- job{data: resp}
			}
		}(i)
	}
	return &s, nil
}

// Handle ...
func (s *Pool) Handle(data interface{}) (interface{}, error) {
	reqCh := make(chan job)
	select {
	case s.ch <- reqCh:
		reqCh <- job{data: data}
		r := <-reqCh
		return r.data, r.err
	default:
		return nil, ErrBufferFull
	}
}

// Len ...
func (s *Pool) Len() int {
	return len(s.ch)
}
