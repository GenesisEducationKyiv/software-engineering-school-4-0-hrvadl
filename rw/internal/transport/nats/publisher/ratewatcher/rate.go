package ratewatcher

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
)

const (
	operation = "rate exchange"
	subject   = "NewExchangeRateFetched"
)

func NewClient(nats *nats.Conn, converter RateSource, log *slog.Logger) *Client {
	return &Client{
		nats:      nats,
		converter: converter,
		log:       log,
	}
}

//go:generate mockgen -destination=./mocks/mock_source.go -package=mocks . RateSource
type RateSource interface {
	Convert(ctx context.Context) (float32, error)
}

type Client struct {
	nats      *nats.Conn
	converter RateSource
	log       *slog.Logger
}

func (c *Client) Convert(ctx context.Context) error {
	_, err := c.converter.Convert(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	if _, err = c.nats.RequestWithContext(ctx, subject, nil); err != nil {
		return fmt.Errorf("%s: failed to send message to NATS: %w", operation, err)
	}

	return nil
}
