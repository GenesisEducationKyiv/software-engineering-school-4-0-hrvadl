package sub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/subscriber"
)

const (
	operation     = "sub subscriber"
	subject       = "subscribers-changed"
	failedSubject = "subscribers-changed-failed"
	stream        = "DebeziumStream"
	consumer      = "sub-consumer"
	deleteEvent   = "delete-subscriber"
	insertEvent   = "add-subscriber"
)

func NewSubscriber(
	conn *nats.Conn,
	sc SubscriberSource,
	log *slog.Logger,
	timeout time.Duration,
) *Subscriber {
	return &Subscriber{
		nats:      conn,
		commander: sc,
		log:       log,
		timeout:   timeout,
	}
}

type SubscriberSaver interface {
	Save(ctx context.Context, sub subscriber.Subscriber) error
}

type SubscriberDeleter interface {
	Delete(ctx context.Context, sub subscriber.Subscriber) error
}

type SubscriberSource interface {
	SubscriberSaver
	SubscriberDeleter
}

type Subscriber struct {
	nats      *nats.Conn
	commander SubscriberSource
	log       *slog.Logger
	timeout   time.Duration
}

func (s *Subscriber) Subscribe() error {
	if _, err := s.nats.Subscribe(subject, s.subscribe); err != nil {
		return fmt.Errorf("%s: failed to consume: %w", operation, err)
	}

	return nil
}

type SubscriberChangedEvent struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Email string `json:"payload"`
}

func (s *Subscriber) subscribe(msg *nats.Msg) {
	var in SubscriberChangedEvent
	if err := json.Unmarshal(msg.Data, &in); err != nil {
		s.log.Error("Failed to parse change event", slog.Any("err", err))
		return
	}

	defer s.ack(msg)
	s.log.Info("Got sub change event from NATS")

	sub := subscriber.Subscriber{Email: in.Email}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	var err error
	switch in.Type {
	case deleteEvent:
		err = s.commander.Delete(ctx, sub)
	case insertEvent:
		err = s.commander.Save(ctx, sub)
	default:
		s.log.Error("Unknown event", slog.Any("type", in.Type))
		return
	}

	if err != nil {
		s.log.Error("Failed to delete/save sub", slog.Any("err", err))
		s.fail(msg.Data)
	}
}

func (s *Subscriber) fail(data []byte) {
	const failTimeout = time.Second * 5
	if _, err := s.nats.Request(failedSubject, data, failTimeout); err != nil {
		s.log.Error("Failed to send fail event", slog.Any("err", err))
	}
}

func (s *Subscriber) ack(msg *nats.Msg) {
	if err := msg.Ack(); err != nil {
		s.log.Error("Failed to send ack", slog.Any("err", err))
	}
}
