package sub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go/jetstream"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/subscriber"
)

const (
	operation = "sub subscriber"
	subject   = "subscribers-changed"
	stream    = "DebeziumStream"
	consumer  = "sub-consumer"
)

func NewSubscriber(
	js jetstream.JetStream,
	sc SubscriberCommander,
	log *slog.Logger,
	timeout time.Duration,
) *Subscriber {
	return &Subscriber{
		stream:    js,
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

type SubscriberCommander interface {
	SubscriberSaver
	SubscriberDeleter
}

type Subscriber struct {
	stream    jetstream.JetStream
	commander SubscriberCommander
	log       *slog.Logger
	timeout   time.Duration
}

func (s *Subscriber) Subscribe() error {
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
	s.log.Info(
		"Got message from NATS",
		slog.Any("msg", string(msg.Data())),
	)

	var in SubscriberChangedEvent
	if err := json.Unmarshal(msg.Data(), &in); err != nil {
		s.log.Error("Failed to parse change event", slog.Any("err", err))
		return
	}

	s.ack(msg)
	sub := subscriber.Subscriber{Email: in.Email}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	var err error
	if in.Deleted {
		err = s.commander.Save(ctx, sub)
	} else {
		err = s.commander.Delete(ctx, sub)
	}

	if err != nil {
		s.log.Error("Failed to delete/save rate", slog.Any("err", err))
	}
}

func (s *Subscriber) ack(msg jetstream.Msg) {
	if err := msg.Ack(); err != nil {
		s.log.Error("Failed to send ack", slog.Any("err", err))
	}
}
