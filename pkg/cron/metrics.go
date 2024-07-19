package cron

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	EventProcessed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "event_sent_total",
		Help: "The total number of sent events",
	}, []string{"status", "event"})

	EventTime = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "event_sent_seconds",
		Help: "The total time of events",
	}, []string{"event"})
)

const (
	statusFailed = "failed"
	statusOK     = "ok"
)

func GetMetrics() []prometheus.Collector {
	return []prometheus.Collector{EventProcessed, EventTime}
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

	EventTime.WithLabelValues(md.event).Observe(time.Since(now).Seconds())
	EventProcessed.With(prometheus.Labels{"status": status, "event": md.event}).Inc()

	return err
}
