package queue

import (
	"hash/fnv"
	"sync"
)

// Item ...
type Item interface {
	String() string
}

// Interface ...
type Interface interface {
	Do(interface{})
}

type queue struct {
	mtx sync.Mutex
	q   map[string]Item

	condMtx sync.Mutex
	cond    *sync.Cond
}

func newQueue() *queue {
	s := queue{q: make(map[string]Item)}
	s.cond = sync.NewCond(&s.condMtx)
	return &s
}

func (s *queue) len() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return len(s.q)
}

func (s *queue) push(data Item) {
	s.mtx.Lock()
	s.q[data.String()] = data
	s.mtx.Unlock()

	s.condMtx.Lock()
	s.cond.Signal()
	s.condMtx.Unlock()
}

func (s *queue) pop() Item {
	s.mtx.Lock()
	if len(s.q) == 0 {
		s.mtx.Unlock()
		s.condMtx.Lock()
		s.cond.Wait()
		s.condMtx.Unlock()
		s.mtx.Lock()
	}

	defer s.mtx.Unlock()
	for k, v := range s.q {
		delete(s.q, k)
		return v
	}
	return nil
}

// Queue ...
type Queue struct {
	maxWorker int
	q         []*queue
}

// NewQueue ...
func NewQueue(maxWorker int, req Interface) *Queue {
	if maxWorker == 0 {
		return nil
	}
	s := Queue{maxWorker: maxWorker}

	for i := 0; i < maxWorker; i++ {
		s.q = append(s.q, newQueue())

		go func(i int) {
			for {
				if r := s.q[i].pop(); r != nil {
					req.Do(r)
				}
			}
		}(i)
	}

	return &s
}

// Handle ...
func (s *Queue) Handle(data Item) {
	h := fnv.New64a()
	h.Write([]byte(data.String()))
	i := h.Sum64() % uint64(s.maxWorker)
	s.q[i].push(data)
}

// Len ...
func (s *Queue) Len() int {
	len := 0
	for i := 0; i < s.maxWorker; i++ {
		len += s.q[i].len()
	}
	return len
}
