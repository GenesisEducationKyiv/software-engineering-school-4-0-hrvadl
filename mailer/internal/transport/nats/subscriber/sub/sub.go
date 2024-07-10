package sub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/subscriber"
)

const (
	operation     = "sub subscriber"
	subject       = "subscribers-changed"
	failedSubject = "subscribers-changed-failed"
	stream        = "DebeziumStream"
	consumer      = "sub-consumer"
)

func NewSubscriber(
	conn *nats.Conn,
	sc SubscriberSource,
	log *slog.Logger,
	timeout time.Duration,
) *Subscriber {
	return &Subscriber{
		conn:      conn,
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
	conn      *nats.Conn
	stream    jetstream.JetStream
	commander SubscriberSource
	log       *slog.Logger
	timeout   time.Duration
}

func (s *Subscriber) Subscribe() error {
	var err error
	if s.stream, err = jetstream.New(s.conn); err != nil {
		return fmt.Errorf("%s: failed to create jetstream: %w", operation, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	stream, err := s.stream.Stream(ctx, stream)
	if err != nil {
		return fmt.Errorf("%s: failed to subscribe to jetstream: %w", operation, err)
	}

	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:          consumer,
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		FilterSubject: subject,
	})
	if err != nil {
		return fmt.Errorf("%s: failed to create a consumer: %w", operation, err)
	}

	if _, err = cons.Consume(s.subscribe); err != nil {
		return fmt.Errorf("%s: failed to consume: %w", operation, err)
	}

	return nil
}

type SubscriberChangedEvent struct {
	Email   string `json:"email"`
	Deleted bool   `json:"__deleted,string"`
}

func (s *Subscriber) subscribe(msg jetstream.Msg) {
	var in SubscriberChangedEvent
	if err := json.Unmarshal(msg.Data(), &in); err != nil {
		s.log.Error("Failed to parse change event", slog.Any("err", err))
		return
	}

	defer s.ack(msg)
	s.log.Info("Got sub change event from NATS", slog.Bool("deleted", in.Deleted))

	sub := subscriber.Subscriber{Email: in.Email}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	var err error
	if in.Deleted {
		err = s.commander.Delete(ctx, sub)
	} else {
		err = s.commander.Save(ctx, sub)
	}

	if err != nil {
		s.log.Error("Failed to delete/save sub", slog.Any("err", err))
		s.fail(msg.Data())
	}
}

func (s *Subscriber) fail(data []byte) {
	const failTimeout = time.Second * 5
	if _, err := s.conn.Request(failedSubject, data, failTimeout); err != nil {
		s.log.Error("Failed to send fail event", slog.Any("err", err))
	}
}

func (s *Subscriber) ack(msg jetstream.Msg) {
	if err := msg.Ack(); err != nil {
		s.log.Error("Failed to send ack", slog.Any("err", err))
	}
}
