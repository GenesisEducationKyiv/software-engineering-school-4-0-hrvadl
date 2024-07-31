package mailer

import (
	"time"

	"github.com/nats-io/nats.go"
)

func NewAdapter(nats *nats.Conn, timeout time.Duration) *Adapter {
	return &Adapter{
		nats:    nats,
		timeout: timeout,
	}
}

type Adapter struct {
	nats    *nats.Conn
	timeout time.Duration
}

func (a *Adapter) Publish(sub string, data []byte) error {
	_, err := a.nats.Request(sub, data, a.timeout)
	return err
}
