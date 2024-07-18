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

func NewServer(addr string) *Prom {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: readTimeout,
	}
	return &Prom{
		srv: srv,
	}
}

type Prom struct {
	srv *http.Server
}

func (p *Prom) Register(c ...prometheus.Collector) error {
	for _, cc := range c {
		if err := prometheus.Register(cc); err != nil {
			return fmt.Errorf("%s: failed to register collector: %w", operation, err)
		}
	}
	return nil
}

func (p *Prom) Start() error {
	if err := p.srv.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: failed to serve metrics: %w", operation, err)
	}
	return nil
}

func (p *Prom) Stop() error {
	return p.srv.Close()
}
