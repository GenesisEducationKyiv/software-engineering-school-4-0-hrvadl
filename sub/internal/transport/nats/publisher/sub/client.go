package sub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
)

const (
	operation = "nats subscription"
	subject   = "subscribers-changed"
)

func NewPublisher(conn *nats.Conn) *EventPublisher {
	return &EventPublisher{
		nats: conn,
	}
}

type EventPublisher struct {
	nats *nats.Conn
}

type SubscriberChangedEvent struct {
	Email   string `json:"email"`
	Deleted bool   `json:"__deleted,string"`
}

func (s *EventPublisher) Publish(ctx context.Context, event SubscriberChangedEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		slog.Info("Failed to marshall payload")
		return fmt.Errorf("%s: failed to marshall payload: %w", operation, err)
	}

	if _, err := s.nats.RequestWithContext(ctx, subject, bytes); err != nil {
		return fmt.Errorf("%s: failed to subscribe to subject: %w", operation, err)
	}

	return nil
}
