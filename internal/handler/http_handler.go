package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/balugcath/william/pkg/metric"
	"github.com/balugcath/william/pkg/queue"
	"github.com/balugcath/william/pkg/types"
	log "github.com/sirupsen/logrus"
)

// HTTPHandler ...
type HTTPHandler struct {
	types.Config
	queue  queue.Interface
	metric metric.Interface
}

// NewHTTPHandler ...
func NewHTTPHandler(c types.Config, q queue.Interface, m metric.Interface) *HTTPHandler {
	s := HTTPHandler{Config: c, queue: q, metric: m}
	s.metric.Register(metric.CounterVec, types.RadAcctMetricName,
		types.RadAcctMetricHelp, []string{"node", "type", "res"}...)
	return &s
}

// Start ...
func (s *HTTPHandler) Start() {
	http.HandleFunc(s.NodeCheckURL, s.nodeChk)
	http.HandleFunc(s.RadiusHTTPAcctURL, s.radAcct)

	go func() {
		if err := http.ListenAndServe(s.RadHTTPListen, nil); err != nil {
			log.Fatal(err)
		}
	}()
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
	s.queue.Push(pckt)
	s.metric.Add(types.RadAcctMetricName, []interface{}{s.NodeName, "received", "ok", float64(1)}...)
	w.WriteHeader(http.StatusNoContent)
	log.Debugf("http handler receive: %+v", pckt)
}

func (s *HTTPHandler) nodeChk(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
	w.WriteHeader(http.StatusNoContent)
}
