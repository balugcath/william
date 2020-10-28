package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/balugcath/william/pkg/metric"
	"github.com/balugcath/william/pkg/queue"
	"github.com/balugcath/william/pkg/types"
	log "github.com/sirupsen/logrus"
)

type radAcctHandler interface {
	Handle(queue.Item)
}

type radAuthHandler interface {
	Handle(interface{}) (interface{}, error)
}

// HTTPHandler ...
type HTTPHandler struct {
	types.Config
	radAcctHandler
	radAuthHandler
	metric metric.Interface
}

// NewHTTPHandler ...
func NewHTTPHandler(c types.Config, acct radAcctHandler, auth radAuthHandler, m metric.Interface) *HTTPHandler {
	s := HTTPHandler{Config: c, radAcctHandler: acct, radAuthHandler: auth, metric: m}
	s.metric.Register(metric.CounterVec, types.RadAcctMetricName,
		types.RadAcctMetricHelp, []string{"node", "type", "res"}...)
	return &s
}

// Start ...
func (s *HTTPHandler) Start() {
	http.HandleFunc(s.NodeCheckURL, s.nodeChk)
	http.HandleFunc(s.RadiusHTTPAcctURL, s.radAcct)
	http.HandleFunc(s.RadiusHTTPAuthURL, s.radAuth)

	go func() {
		if err := http.ListenAndServe(s.RadHTTPListen, nil); err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *HTTPHandler) radAuth(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Header.Get(s.RadHTTPHeader) != s.RadHTTPToken {
		http.NotFound(w, r)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("%w %s", types.ErrHTTPHandler, err).Error(), http.StatusBadRequest)
		s.metric.Add(types.RadAuthMetricName, []interface{}{s.NodeName, "received", "err", float64(1)}...)
		log.Errorf("%s %s", types.ErrHTTPHandler, err)
		return
	}
	log.Debugf("http handler rad auth receive: %+v", string(body))

	s.metric.Add(types.RadAuthMetricName, []interface{}{s.NodeName, "received", "ok", float64(1)}...)
	resp, err := s.radAuthHandler.Handle(string(body))
	if err != nil {
		if errors.Is(err, types.ErrRadAuthReject) {
			log.Debugf("http handler rad auth reject: %+v", string(body))
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(resp.(string)))
			s.metric.Add(types.RadAuthMetricName, []interface{}{s.NodeName, "reply", "reject", float64(1)}...)
			return
		}
		log.Debugf("http handler rad auth error: %+v", string(body))
		s.metric.Add(types.RadAuthMetricName, []interface{}{s.NodeName, "reply", "error", float64(1)}...)
		http.Error(w, fmt.Errorf("%w %s", types.ErrHTTPHandler, err).Error(), http.StatusUnauthorized)
		log.Errorf("%s %s", types.ErrHTTPHandler, err)
		return
	}
	log.Debugf("http handler rad auth ok: %+v", string(body))
	s.metric.Add(types.RadAuthMetricName, []interface{}{s.NodeName, "reply", "ok", float64(1)}...)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp.(string)))
}

func (s *HTTPHandler) radAcct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Header.Get(s.RadHTTPHeader) != s.RadHTTPToken {
		http.NotFound(w, r)
		return
	}

	pckt := types.RadiusAccounting{}
	if err := json.NewDecoder(r.Body).Decode(&pckt); err != nil {
		http.Error(w, fmt.Errorf("%w %s", types.ErrHTTPHandler, err).Error(), http.StatusBadRequest)
		s.metric.Add(types.RadAcctMetricName, []interface{}{s.NodeName, "received", "err", float64(1)}...)
		log.Errorf("%s %s", types.ErrHTTPHandler, err)
		return
	}
	s.radAcctHandler.Handle(pckt)
	s.metric.Add(types.RadAcctMetricName, []interface{}{s.NodeName, "received", "ok", float64(1)}...)
	w.WriteHeader(http.StatusNoContent)
	log.Debugf("http handler rad acct receive: %+v", pckt)
}

func (s *HTTPHandler) nodeChk(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
	w.WriteHeader(http.StatusNoContent)
}
