package sub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/event"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/subscriber"
)

const (
	operation   = "nats subscription"
	subject     = "subscribers-changed-failed"
	deleteEvent = "delete-subscriber"
	insertEvent = "add-subscriber"
)

func NewSubscriber(conn *nats.Conn, compensator Compensator, log *slog.Logger) *EventSubscriber {
	return &EventSubscriber{
		nats:        conn,
		log:         log,
		compensator: compensator,
	}
}

type Subscriber interface {
	Subscribe(ctx context.Context, sub subscriber.Subscriber) (int64, error)
}

type Unsubscriber interface {
	Unsubscribe(ctx context.Context, sub subscriber.Subscriber) error
}

type Compensator interface {
	Subscriber
	Unsubscriber
}

type EventSubscriber struct {
	nats        *nats.Conn
	log         *slog.Logger
	compensator Compensator
}

func (s *EventSubscriber) Subscribe() error {
	_, err := s.nats.Subscribe(subject, s.consume)
	if err != nil {
		return fmt.Errorf("%s: failed to subscribe to subject: %w", operation, err)
	}
	return nil
}

type SubscriberChangedEvent struct {
	ID    int        `json:"id"`
	Type  event.Type `json:"type"`
	Email string     `json:"payload"`
}

func (s *EventSubscriber) consume(msg *nats.Msg) {
	const compensateTimeout = time.Second * 5

	var in SubscriberChangedEvent
	if err := json.Unmarshal(msg.Data, &in); err != nil {
		s.log.Error("Failed to parse change event", slog.Any("err", err))
		return
	}

	s.log.Info("Got event from event bus")

	ctx, cancel := context.WithTimeout(context.Background(), compensateTimeout)
	defer cancel()

	var (
		err error
		sub = subscriber.Subscriber{Email: in.Email}
	)

	switch in.Type {
	case event.Deleted:
		_, err = s.compensator.Subscribe(ctx, sub)
	case event.Added:
		err = s.compensator.Unsubscribe(ctx, sub)
	default:
		s.log.Error("Unknown event", slog.Any("type", in.Type))
		s.ack(msg)
		return
	}

	if err != nil {
		s.log.Error("Failed to compensate transaction", slog.Any("err", err))
		s.nack(msg)
		return
	}

	s.ack(msg)
}

func (s *EventSubscriber) ack(msg *nats.Msg) {
	if err := msg.Ack(); err != nil {
		s.log.Error("Failed to send ack", slog.Any("err", err))
	}
}

func (s *EventSubscriber) nack(msg *nats.Msg) {
	if err := msg.Nak(); err != nil {
		s.log.Error("Failed to send nack", slog.Any("err", err))
	}
}
