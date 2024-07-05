package mailer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

const (
	subjectSub = "new-sub.converter.subscribers"
	stream     = "DebeziumStream"
)

func NewSub(js jetstream.JetStream, log *slog.Logger) *ServerSub {
	return &ServerSub{
		stream: js,
		log:    log,
	}
}

type ServerSub struct {
	stream jetstream.JetStream
	log    *slog.Logger
}

func (s *ServerSub) Subscribe() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := s.stream.Stream(ctx, stream)
	if err != nil {
		return fmt.Errorf("%s: failed to subscribe to jetstream: %w", operation, err)
	}

	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:          "SubConsumer",
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		FilterSubject: subjectSub,
	})
	if err != nil {
		return fmt.Errorf("%s: failed to create a consumer: %w", operation, err)
	}

	if _, err = cons.Consume(s.subscribe); err != nil {
		return fmt.Errorf("%s: failed to consume: %w", operation, err)
	}

	return nil
}

func (s *ServerSub) subscribe(msg jetstream.Msg) {
	s.log.Info(
		"Got message from NATS",
		slog.Any("msg", string(msg.Data())),
	)

	s.ack(msg)
}

func (s *ServerSub) ack(msg jetstream.Msg) {
	if err := msg.Ack(); err != nil {
		s.log.Error("Failed to send ack", slog.Any("err", err))
	}
}
