package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/balugcath/william/pkg/metric"
	"github.com/balugcath/william/pkg/queue"
	"github.com/balugcath/william/pkg/types"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type userIDHandler interface {
	Handle(queue.Item)
}

// SQLListenHandler ...
type SQLListenHandler struct {
	types.Config
	userIDHandler
	metric metric.Interface
}

const (
	minReconn = time.Minute
	maxReconn = time.Minute * 10
)

// NewSQLListenHandler ...
func NewSQLListenHandler(c types.Config, userID userIDHandler, m metric.Interface) *SQLListenHandler {
	s := SQLListenHandler{Config: c, userIDHandler: userID, metric: m}
	s.metric.Register(metric.CounterVec, types.UserIDMetricName, types.UserIDMetricHelp,
		[]string{"node", "type", "res"}...)
	return &s
}

// Start ...
func (s SQLListenHandler) Start() error {
	l, err := s.listen()
	if err != nil {
		return fmt.Errorf("%w %s", types.ErrSQLListen, err)
	}

	go func() {
		defer l.Close()

		for {
			p, err := s.get(l)
			if err != nil {
				if err == types.ErrSQLListenTimeout {
					continue
				}
				s.metric.Add(types.UserIDMetricName, []interface{}{s.NodeName, "received", "err", float64(1)}...)
				log.Error(err)
				continue
			}
			log.Debugf("sql listen handler receive: %+v", p)
			s.userIDHandler.Handle(p)
			s.metric.Add(types.UserIDMetricName, []interface{}{s.NodeName, "received", "ok", float64(1)}...)
		}

	}()
	return nil
}

func (s SQLListenHandler) listen() (*pq.Listener, error) {
	l := pq.NewListener(s.DBURI, minReconn, maxReconn, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Fatal(fmt.Errorf("%w %s", types.ErrSQLListen, err))
		}
	})
	return l, l.Listen(s.SQLListenChan)
}

func (s SQLListenHandler) get(l *pq.Listener) (types.UserID, error) {
	select {
	case recv := <-l.NotificationChannel():
		id, err := strconv.Atoi(recv.Extra)
		if err != nil {
			return types.UserID{}, fmt.Errorf("%w %s", types.ErrSQLListen, err)
		}
		return types.UserID{UserID: id}, nil
	case <-time.After(time.Second * 120):
		go l.Ping()
		return types.UserID{}, types.ErrSQLListenTimeout
	}
}
