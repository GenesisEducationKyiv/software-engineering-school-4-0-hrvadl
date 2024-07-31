package sub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/sub"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/subscriber"
)

const (
	operation     = "sub subscriber"
	subject       = "subscribers-changed"
	failedSubject = "subscribers-changed-failed"
	stream        = "DebeziumStream"
	consumer      = "sub-consumer"
	deleteEvent   = "subscriber-deleted"
	insertEvent   = "subscriber-added"
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
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func (s *Subscriber) subscribe(msg *nats.Msg) {
	var in SubscriberChangedEvent
	if err := json.Unmarshal(msg.Data, &in); err != nil {
		s.log.Error(
			"Failed to parse change event",
			slog.Any("err", err),
			slog.String("data", string(msg.Data)),
		)
		return
	}

	defer s.ack(msg)
	s.log.Info("Got sub change event from NATS")

	var ev pb.SubscriptionAddedEvent
	if err := protojson.Unmarshal([]byte(in.Payload), &ev); err != nil {
		s.log.Error("Failed to parse change event", slog.Any("err", err))
		return
	}

	sub := subscriber.Subscriber{Email: ev.GetEmail()}
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
		return
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
