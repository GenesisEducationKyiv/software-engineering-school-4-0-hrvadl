package sub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/subscriber"
)

const (
	operation = "nats subscription"
	subject   = "subscribers-changed-failed"
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
	Email   string `json:"email"`
	Deleted bool   `json:"__deleted,string"`
}

func (s *EventSubscriber) consume(msg *nats.Msg) {
	const compensateTimeout = time.Second * 5

	var in SubscriberChangedEvent
	if err := json.Unmarshal(msg.Data, &in); err != nil {
		s.log.Error("Failed to parse change event", slog.Any("err", err))
		return
	}

	log := s.log.With(slog.Bool("deleted", in.Deleted))
	log.Info("Got event from event bus")

	defer s.ack(msg)
	ctx, cancel := context.WithTimeout(context.Background(), compensateTimeout)
	defer cancel()

	var (
		err error
		sub = subscriber.Subscriber{Email: in.Email}
	)

	if in.Deleted {
		_, err = s.compensator.Subscribe(ctx, sub)
	} else {
		err = s.compensator.Unsubscribe(ctx, sub)
	}

	if err != nil {
		log.Error("Failed to compensate transaction", slog.Any("err", err))
	}
}

func (s *EventSubscriber) ack(msg *nats.Msg) {
	if err := msg.Ack(); err != nil {
		s.log.Error("Failed to send ack", slog.Any("err", err))
	}
}
