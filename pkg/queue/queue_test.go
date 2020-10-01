package queue

// import (
// 	"runtime"
// 	"strconv"
// 	"sync"
// 	"testing"
// )

// type push struct {
// 	mtx sync.Mutex
// 	res []interface{}
// }

// func (s *push) Push(r interface{}) {
// 	s.mtx.Lock()
// 	s.res = append(s.res, r)
// 	s.mtx.Unlock()
// }

// type data struct {
// 	id     int
// 	series int
// }

// func (s data) String() string {
// 	return strconv.Itoa(s.id)
// }

// func Test_queue_add_fetch_1(t *testing.T) {
// 	type args struct {
// 		cntGenItem int
// 		push       *push
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		{
// 			name: "test1",
// 			args: args{
// 				cntGenItem: 8,
// 				push:       &push{res: make([]item, 0)},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		ch := make(chan int, 1)
// 		s := newQueue()
// 		go func() {
// 			tt.args.push.Push(s.fetch())
// 		}()
// 		t.Run(tt.name, func(t *testing.T) {
// 			go func() {
// 				for i := 0; i < tt.args.cntGenItem; i++ {
// 					s.add(data{id: 1, series: i})
// 				}
// 				ch <- 1
// 			}()
// 			<-ch
// 			runtime.Gosched()
// 			tt.args.push.mtx.Lock()
// 			t.Logf("res = %+v", tt.args.push.res)
// 			tt.args.push.mtx.Unlock()
// 		})
// 	}
// }

// func Test_queue_add_fetch_1(t *testing.T) {
// 	type args struct {
// 		cntGenItem int
// 		push       *push
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		{
// 			name: "test1",
// 			args: args{
// 				cntGenItem: 2,
// 				push:       &push{res: make([]item, 0)},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		ch := make(chan int, 2)

// 		t.Run(tt.name, func(t *testing.T) {
// 			s := newQueue()
// 			go func() {
// 				for i := 0; i < tt.args.cntGenItem; i++ {
// 					go func(i int) {
// 						s.add(func() item {
// 							return data{id: i, series: 0}
// 						}())
// 					}(i)
// 					tt.args.push.Push(s.fetch())
// 				}
// 				ch <- 1
// 			}()
// 			go func() {
// 				for i := 0; i < tt.args.cntGenItem; i++ {
// 					go func() {
// 						tt.args.push.Push(s.fetch())
// 					}()
// 					s.add(func() item {
// 						return data{id: i, series: 1}
// 					}())
// 				}
// 				ch <- 1
// 			}()
// 			<-ch
// 			<-ch
// 			tt.args.push.mtx.Lock()
// 			t.Logf("res = %+v", tt.args.push.res)
// 			tt.args.push.mtx.Unlock()
// 		})
// 	}

// // }

// func TestQueue_QueuePush(t *testing.T) {
// 	type args struct {
// 		maxWorker int
// 		push      *push
// 		req       []data
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []data
// 	}{
// 		{
// 			name: "test1",
// 			args: args{
// 				maxWorker: 1,
// 				push:      &push{res: make([]interface{}, 0)},
// 				req:       []data{{id: 1}, {id: 2}, {id: 3}},
// 			},
// 			want: []data{{id: 1}, {id: 2}, {id: 3}},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := NewQueue(tt.args.maxWorker, tt.args.push).Start()

// 			for i := range tt.args.req {
// 				s.Add(tt.args.req[i])
// 			}
// 			runtime.Gosched()

// 			t.Logf("1 %+v", s.q[0].q)
// 			t.Logf("2 %+v", tt.args.push.res)
// 		})
// 	}
// }
