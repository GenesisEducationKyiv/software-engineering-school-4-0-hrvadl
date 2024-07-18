package cron

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	mailsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mail_sent_total",
		Help: "The total number of sent mail events",
	}, []string{"status"})

	mailTime = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "mail_sent_seconds",
		Help: "The total time of mail rate events",
	})
)

const (
	statusFailed = "failed"
	statusOK     = "ok"
)

func NewWithMetrics(doer Doer) *MetricsDecorator {
	return &MetricsDecorator{
		doer: doer,
	}
}

//go:generate mockgen -destination=./mocks/mock_doer.go -package=mocks . Doer
type Doer interface {
	Do() error
}

type MetricsDecorator struct {
	doer Doer
}

func (md *MetricsDecorator) Do() error {
	var (
		now    = time.Now()
		status = statusOK
	)

	err := md.doer.Do()
	if err != nil {
		status = statusFailed
	}

	mailTime.Observe(time.Since(now).Seconds())
	mailsProcessed.With(prometheus.Labels{"status": status}).Inc()

	return err
}
