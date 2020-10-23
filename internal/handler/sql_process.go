package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/balugcath/william/pkg/metric"
	"github.com/balugcath/william/pkg/types"
	log "github.com/sirupsen/logrus"
)

// SQLProcessRadAcct ...
type SQLProcessRadAcct struct {
	types.Config
	db     *sql.DB
	metric metric.Interface
}

// NewSQLProcessRadAcct ...
func NewSQLProcessRadAcct(db *sql.DB, c types.Config, m metric.Interface) *SQLProcessRadAcct {
	s := SQLProcessRadAcct{db: db, Config: c, metric: m}
	s.metric.Register(metric.CounterVec, types.RadAcctMetricName,
		types.RadAcctMetricHelp, []string{"node", "type", "res"}...)
	return &s
}

// Do ...
func (s *SQLProcessRadAcct) Do(r interface{}) {
	b, err := json.Marshal(r)
	if err != nil {
		s.metric.Add(types.RadAcctMetricName, []interface{}{s.NodeName, "processed", "err", float64(1)}...)
		log.Error(fmt.Errorf("%w %s", types.ErrSQLProcess, err))
		return
	}
	_, err = s.db.Exec(s.RadAcctSQLQuery, string(b))
	if err != nil {
		s.metric.Add(types.RadAcctMetricName, []interface{}{s.NodeName, "processed", "err", float64(1)}...)
		log.Error(fmt.Errorf("%w %s", types.ErrSQLProcess, err))
		return
	}
	s.metric.Add(types.RadAcctMetricName, []interface{}{s.NodeName, "processed", "ok", float64(1)}...)
	log.Debugf("sql process rad acct req %+v", string(b))
}

// SQLProcessUserID ...
type SQLProcessUserID struct {
	db *sql.DB
	types.Config
	metric metric.Interface
}

// NewSQLProcessUserID ...
func NewSQLProcessUserID(db *sql.DB, c types.Config, m metric.Interface) *SQLProcessUserID {
	s := SQLProcessUserID{db: db, Config: c, metric: m}
	s.metric.Register(metric.CounterVec, types.UserIDMetricName, types.UserIDMetricHelp,
		[]string{"node", "type", "res"}...)
	return &s
}

// Do ...
func (s *SQLProcessUserID) Do(r interface{}) {
	b, err := json.Marshal(r)
	if err != nil {
		s.metric.Add(types.UserIDMetricName, []interface{}{s.NodeName, "processed", "err", float64(1)}...)
		log.Error(fmt.Errorf("%w %s", types.ErrSQLProcess, err))
		return
	}
	_, err = s.db.Exec(s.UserIDSQLQuery, string(b))
	if err != nil {
		s.metric.Add(types.UserIDMetricName, []interface{}{s.NodeName, "processed", "err", float64(1)}...)
		log.Error(fmt.Errorf("%w %s", types.ErrSQLProcess, err))
		return
	}
	s.metric.Add(types.UserIDMetricName, []interface{}{s.NodeName, "processed", "ok", float64(1)}...)
	log.Debugf("sql process user_id req %+v", string(b))
}

// SQLProcessRadAuth ...
type SQLProcessRadAuth struct {
	db *sql.DB
	types.Config
	metric metric.Interface
}

// NewSQLProcessRadAuth ...
func NewSQLProcessRadAuth(db *sql.DB, c types.Config, m metric.Interface) *SQLProcessRadAuth {
	s := SQLProcessRadAuth{db: db, Config: c, metric: m}
	s.metric.Register(metric.CounterVec, types.RadAuthMetricName, types.RadAuthMetricHelp,
		[]string{"node", "type", "res"}...)
	return &s
}

// Do ...
func (s *SQLProcessRadAuth) Do(r interface{}) (interface{}, error) {
	var (
		resp     string
		isReject bool
	)
	req, ok := r.(string)
	if !ok {
		return nil, fmt.Errorf("%w wrong parameter ", types.ErrSQLProcess)
	}
	err := s.db.QueryRow(s.RadAuthSQLQuery, req).Scan(&resp, &isReject)
	if err != nil {
		s.metric.Add(types.RadAuthMetricName, []interface{}{s.NodeName, "processed", "err", float64(1)}...)
		log.Error(fmt.Errorf("%w %s", types.ErrSQLProcess, err))
		return nil, fmt.Errorf("%w %s", types.ErrSQLProcess, err)
	}
	s.metric.Add(types.RadAuthMetricName, []interface{}{s.NodeName, "processed", "ok", float64(1)}...)
	log.Debugf("sql process radauth req %+v", req)
	if isReject {
		return resp, types.ErrRadAuthReject
	}
	return resp, nil
}
