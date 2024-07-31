package cron

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	eventProcessed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "event_sent_total",
		Help: "The total number of sent events",
	}, []string{"status", "event"})

	eventTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "event_sent_seconds",
		Help:    "The total time of events",
		Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1},
	}, []string{"event"})
)

const (
	statusFailed = "failed"
	statusOK     = "ok"
)

func GetMetrics() []prometheus.Collector {
	return []prometheus.Collector{eventProcessed, eventTime}
}

func NewWithMetrics(doer Doer, event string) *MetricsDecorator {
	return &MetricsDecorator{
		doer:  doer,
		event: event,
	}
}

type MetricsDecorator struct {
	doer  Doer
	event string
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

	eventTime.WithLabelValues(md.event).Observe(time.Since(now).Seconds())
	eventProcessed.With(prometheus.Labels{"status": status, "event": md.event}).Inc()

	return err
}
