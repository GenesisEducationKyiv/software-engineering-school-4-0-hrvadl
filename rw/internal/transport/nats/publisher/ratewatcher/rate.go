package ratewatcher

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v3/mailer"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const (
	operation = "rate exchange"
	subject   = "rate-fetched"
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
	rate, err := c.converter.Convert(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	bytes, err := proto.Marshal(&pb.ExchangeFetchedEvent{
		From:      "USD",
		To:        "UAH",
		Rate:      rate,
		EventID:   "1",
		EventType: operation,
	})
	if err != nil {
		return fmt.Errorf("%s: failed to marshall proto: %w", operation, err)
	}

	if _, err = c.nats.RequestWithContext(ctx, subject, bytes); err != nil {
		return fmt.Errorf("%s: failed to send message to NATS: %w", operation, err)
	}

	return nil
}
