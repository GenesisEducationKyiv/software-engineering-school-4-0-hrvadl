package rate

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v3/mailer"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/rate"
)

const (
	operation = "rate subscriber"
	subject   = "rate-fetched"
	eventType = subject
)

func NewSubscriber(
	conn *nats.Conn,
	s RateReplacer,
	log *slog.Logger,
	timeout time.Duration,
) *Subscriber {
	return &Subscriber{
		conn:     conn,
		replacer: s,
		log:      log,
		timeout:  timeout,
	}
}

type RateReplacer interface {
	Replace(ctx context.Context, r rate.Exchange) error
}

type Subscriber struct {
	conn     *nats.Conn
	replacer RateReplacer
	timeout  time.Duration
	log      *slog.Logger
}

func (s *Subscriber) Subscribe() error {
	_, err := s.conn.Subscribe(subject, s.subscribe)
	if err != nil {
		return fmt.Errorf("%s: failed to subscribe to nats: %w", operation, err)
	}
	return nil
}

func (s *Subscriber) subscribe(msg *nats.Msg) {
	var in pb.ExchangeFetchedEvent
	if err := proto.Unmarshal(msg.Data, &in); err != nil {
		s.log.Error("Failed to parse mail", slog.Any("err", err))
		return
	}

	log := s.log.With(
		slog.String("id", in.GetEventID()),
		slog.String("type", in.GetEventType()),
	)

	if in.GetEventType() != eventType {
		log.Info("Discarding unknown event...")
		s.nack(msg)
		return
	}

	log.Info("Got message from nats")
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	rate := rate.Exchange{
		From: in.GetFrom(),
		To:   in.GetTo(),
		Rate: in.GetRate(),
	}

	if err := s.replacer.Replace(ctx, rate); err != nil {
		log.Error("Failed to replace rate", slog.Any("err", err))
		s.nack(msg)
		return
	}

	s.ack(msg)
	log.Info("Successfully replaced rate")
}

func (s *Subscriber) ack(msg *nats.Msg) {
	if err := msg.Ack(); err != nil {
		s.log.Error("Failed to send ack", slog.Any("err", err))
	}
}

func (s *Subscriber) nack(msg *nats.Msg) {
	if err := msg.Nak(); err != nil {
		s.log.Error("Failed to send nack", slog.Any("err", err))
	}
}
