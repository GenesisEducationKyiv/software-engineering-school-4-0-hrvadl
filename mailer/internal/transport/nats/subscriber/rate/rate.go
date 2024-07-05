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

	defer s.ack(msg)
	s.log.Info(
		"Got message from NATS",
		slog.String("id", in.GetEventID()),
		slog.String("type", in.GetEventType()),
	)

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	rate := rate.Exchange{
		From: in.GetFrom(),
		To:   in.GetTo(),
		Rate: in.GetRate(),
	}

	if err := s.replacer.Replace(ctx, rate); err != nil {
		s.log.Error("Failed to replace rate", slog.Any("err", err))
	}
}

func (s *Subscriber) ack(msg *nats.Msg) {
	if err := msg.Ack(); err != nil {
		s.log.Error("Failed to send ack", slog.Any("err", err))
	}
}
