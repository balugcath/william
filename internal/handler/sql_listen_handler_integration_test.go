// +built integration_test

package handler

// import (
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/balugcath/william/pkg/types"
// )

// func TestSQLListenHandler_Start(t *testing.T) {
// 	type args struct {
// 		dburi      string
// 		listenChan string
// 		req        int
// 		q          *queueMock
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want item
// 	}{
// 		{
// 			name: "test1",
// 			args: args{
// 				dburi:      "user=postgres password=123 dbname=postgres host=192.168.1.31 port=5432 sslmode=disable",
// 				listenChan: "test",
// 				q:          &queueMock{},
// 				req:        123,
// 			},
// 			want: types.UserID{UserID: 123},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			SQLListenHandler{}.Start(tt.args.dburi, tt.args.listenChan, tt.args.q)

// 			db, err := sql.Open("postgres", tt.args.dburi)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			defer db.Close()

// 			_, err = db.Exec(fmt.Sprintf("notify %s, '%d'", tt.args.listenChan, tt.args.req))
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			time.Sleep(time.Second * 2)

// 			if tt.want != tt.args.q.res {
// 				t.Errorf("want %+v got %+v", tt.want, tt.args.q.res)
// 			}

// 		})
// 	}
// }

// func TestSQLListenHandler_listen_get(t *testing.T) {
// 	type args struct {
// 		dburi      string
// 		listenChan string
// 		req        string
// 	}
// 	tests := []struct {
// 		name     string
// 		args     args
// 		want     item
// 		wantErr1 error
// 		wantErr2 error
// 		wantErr3 error
// 	}{
// 		{
// 			name: "test1",
// 			args: args{
// 				dburi:      "user=postgres password=123 dbname=postgres host=192.168.1.31 port=5432 sslmode=disable",
// 				listenChan: "test",
// 				req:        "123",
// 			},
// 			want: types.UserID{UserID: 123},
// 		},
// 		{
// 			name: "test2",
// 			args: args{
// 				dburi:      "user=postgres password=123 dbname=postgres host=192.168.1.31 port=5432 sslmode=disable",
// 				listenChan: "test",
// 				req:        "abc",
// 			},
// 			wantErr2: types.ErrSQLListen,
// 			want:     types.UserID{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			l, err := SQLListenHandler{}.listen(tt.args.dburi, tt.args.listenChan)
// 			if err != tt.wantErr1 {
// 				t.Errorf("want %+v got %+v", tt.wantErr1, err)
// 			}

// 			ch := make(chan error)
// 			go func() {
// 				db, err := sql.Open("postgres", tt.args.dburi)
// 				if err != nil {
// 					ch <- err
// 					return
// 				}
// 				defer db.Close()

// 				_, err = db.Exec(fmt.Sprintf("notify %s, '%s'", tt.args.listenChan, tt.args.req))
// 				if err != nil {
// 					ch <- err
// 					return
// 				}
// 				ch <- nil
// 			}()

// 			err = <-ch
// 			if err != nil && tt.wantErr3 == nil {
// 				t.Errorf("%w unexpected error", err)
// 			}

// 			p, err := SQLListenHandler{}.get(l)
// 			if err != nil && !errors.Is(err, types.ErrSQLListen) {
// 				t.Errorf("want %+v got %+v", tt.wantErr1, err)
// 			}

// 			if p != tt.want {
// 				t.Errorf("want %+v got %+v", tt.want, p)
// 			}

// 		})
// 	}
// }
