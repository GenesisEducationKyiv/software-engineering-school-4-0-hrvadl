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

func NewEngine(addr string) *Engine {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: readTimeout,
	}

	return &Engine{srv: srv}
}

type Engine struct {
	srv *http.Server
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
