package cron

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var mailsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "mail_sent_total",
	Help: "The total number of sent mail events",
}, []string{"status"})

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
	if err := md.doer.Do(); err != nil {
		mailsProcessed.With(prometheus.Labels{"status": statusFailed}).Inc()
		return err
	}

	mailsProcessed.With(prometheus.Labels{"status": statusOK}).Inc()
	return nil
}
