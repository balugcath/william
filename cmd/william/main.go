package main

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/balugcath/william/internal/handler"
	"github.com/balugcath/william/internal/runtime_metric"
	"github.com/balugcath/william/pkg/metric"
	"github.com/balugcath/william/pkg/queue"
	"github.com/balugcath/william/pkg/types"
	"github.com/balugcath/william/pkg/work_pool"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func main() {
	c := types.Config{}

	if err := godotenv.Load(); err != nil {
		log.Info("error loading .env file")
	}

	if err := envconfig.Process("bill", &c); err != nil {
		log.Fatal(err.Error())
	}

	if c.DebugLevel {
		log.SetLevel(log.DebugLevel)
		log.Debugf("config: %+v", c)
	}

	db, err := sql.Open("postgres", c.DBURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	m := metric.NewMetric(c.PrometheusListen, c.PrometheusPath)
	go func() {
		log.Fatal(m.Start())
	}()

	qRadAcct := queue.NewQueue(c.AcctCntWorker,
		handler.NewSQLProcessRadAcct(db, c, m),
	)

	qUserID := queue.NewQueue(c.UserIDCntWorker,
		handler.NewSQLProcessUserID(db, c, m),
	)

	poolRadAuth, _ := workpool.NewPool(c.UserIDCntWorker, c.AuthBuffLen,
		handler.NewSQLProcessRadAuth(db, c, m),
	)

	handler.NewHTTPHandler(c, qRadAcct, poolRadAuth, m).Start()
	handler.NewSQLListenHandler(c, qUserID, m).Start()

	rtmetric.NewRTMetric(c, qRadAcct, m, "radius_acct")
	rtmetric.NewRTMetric(c, poolRadAuth, m, "radius_auth")
	rtmetric.NewRTMetric(c, qUserID, m, "user_id")

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGUSR1)
	for {
		switch <-sigs {
		case syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGUSR1:
			if c.DebugLevel {
				log.SetLevel(log.ErrorLevel)
			} else {
				log.SetLevel(log.DebugLevel)
			}
			c.DebugLevel = !c.DebugLevel
		}
	}
}
