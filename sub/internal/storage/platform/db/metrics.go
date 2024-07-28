package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsProcessed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_sent_total",
		Help: "The total number of sent requests",
	}, []string{"status", "table"})

	requestTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_sent_seconds",
		Help:    "The total time of requests",
		Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1},
	}, []string{"table"})
)

const (
	statusFailed = "failed"
	statusOK     = "ok"
)

func NewWithMetrics(db DataSource) *MetricsDecorator {
	return &MetricsDecorator{
		db: db,
	}
}

type DataSource interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	NamedExecContext(ctx context.Context, query string, args any) (sql.Result, error)
}

type MetricsDecorator struct {
	table string
	db    DataSource
}

func (d MetricsDecorator) WithTableName(tablename string) *MetricsDecorator {
	d.table = tablename
	return &d
}

func (d *MetricsDecorator) GetMetrics() []prometheus.Collector {
	return []prometheus.Collector{requestTime, requestsProcessed}
}

func (d *MetricsDecorator) ExecContext(
	ctx context.Context,
	query string,
	args ...any,
) (sql.Result, error) {
	now := time.Now()
	res, err := d.db.ExecContext(ctx, query, args...)
	d.captureTimeMetric(now)
	d.captureCountMetric(err)
	return res, err
}

func (d *MetricsDecorator) GetContext(
	ctx context.Context,
	dest any,
	query string,
	args ...any,
) error {
	now := time.Now()
	err := d.db.GetContext(ctx, dest, query, args...)
	d.captureTimeMetric(now)
	d.captureCountMetric(err)
	return err
}

func (d *MetricsDecorator) SelectContext(
	ctx context.Context,
	dest any,
	query string,
	args ...any,
) error {
	return d.db.SelectContext(ctx, dest, query, args...)
}

func (d *MetricsDecorator) NamedExecContext(
	ctx context.Context,
	query string,
	arg any,
) (sql.Result, error) {
	now := time.Now()
	res, err := d.db.NamedExecContext(ctx, query, arg)
	d.captureTimeMetric(now)
	d.captureCountMetric(err)
	return res, err
}

func (d *MetricsDecorator) captureTimeMetric(since time.Time) {
	requestTime.WithLabelValues(d.table).Observe(time.Since(since).Seconds())
}

func (d *MetricsDecorator) captureCountMetric(err error) {
	requestsProcessed.WithLabelValues(d.getStatusLabelBasedOnErr(err), d.table).Inc()
}

func (d *MetricsDecorator) getStatusLabelBasedOnErr(err error) string {
	if err != nil {
		return statusFailed
	}
	return statusOK
}
