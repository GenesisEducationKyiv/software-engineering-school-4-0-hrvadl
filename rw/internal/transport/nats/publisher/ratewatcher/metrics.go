package ratewatcher

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	statusFailed = "failed"
	statusOK     = "ok"
)

var (
	rateProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "rate_sent_total",
		Help: "The total number of sent rate events",
	}, []string{"status"})

	rateTime = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "rate_sent_seconds",
		Help: "The total time of sent rate events",
	})
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

	rateTime.Observe(time.Since(now).Seconds())
	rateProcessed.With(prometheus.Labels{"status": status}).Inc()

	return err
}
