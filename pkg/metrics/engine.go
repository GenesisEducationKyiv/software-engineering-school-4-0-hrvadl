package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	readTimeout = time.Second * 5
	operation   = "metrics"
)

var requestTimeDefault = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "request_seconds",
	Help:    "The total time of seconds spent waiting for request",
	Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1},
}, []string{"type", "path", "status"})

func NewEngine(addr string) (*Engine, error) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: readTimeout,
	}

	engine := &Engine{srv: srv}
	if err := engine.registerDefaultMetrics(); err != nil {
		return nil, err
	}

	return engine, nil
}

type Labels = map[string]string

type Engine struct {
	srv *http.Server
}

// CollectRequestTimeWithLabels is experimental approach (alternative to decorators) proposed by my mentor.
// I'm just playing with it to understant which approach suites better to
// use in my application.
func (p *Engine) CollectRequestTimeWithLabels(l Labels, t time.Duration) {
	requestTimeDefault.With(l).Observe(t.Seconds())
}

func (p *Engine) Register(c ...prometheus.Collector) error {
	for _, cc := range c {
		if err := prometheus.Register(cc); err != nil {
			return fmt.Errorf("%s: failed to register collector: %w", operation, err)
		}
	}
	return nil
}

func (p *Engine) Start() error {
	if err := p.srv.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: failed to serve metrics: %w", operation, err)
	}
	return nil
}

func (p *Engine) Stop() error {
	return p.srv.Close()
}

func (p *Engine) registerDefaultMetrics() error {
	if err := prometheus.Register(requestTimeDefault); err != nil {
		return fmt.Errorf("%s: failed to register collector: %w", operation, err)
	}
	return nil
}
