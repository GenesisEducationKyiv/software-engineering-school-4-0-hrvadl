package sub

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/event"
)

const timeout = time.Second * 5

func NewOutboxer(pub Publisher, eg EventGetter, ed EventDeletter, log *slog.Logger) *Outboxer {
	return &Outboxer{
		pub:           pub,
		eventgetter:   eg,
		eventdeletter: ed,
		log:           log,
	}
}

type EventGetter interface {
	GetAll(ctx context.Context) ([]event.Event, error)
}

type EventDeletter interface {
	DeleteByID(ctx context.Context, id int) error
}

type Publisher interface {
	Publish(ctx context.Context, event SubscriberChangedEvent) error
}

type Outboxer struct {
	pub           Publisher
	eventgetter   EventGetter
	eventdeletter EventDeletter
	log           *slog.Logger
}

func (o *Outboxer) Do() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	events, err := o.eventgetter.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get events: %w", err)
	}

	g := new(errgroup.Group)
	for _, e := range events {
		g.Go(func() error {
			return o.publish(e)
		})
	}

	return g.Wait()
}

func (o *Outboxer) publish(e event.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	o.log.Info("Running outbox...")
	err := o.pub.Publish(ctx, SubscriberChangedEvent{
		Email:   e.Payload,
		Deleted: e.Type == event.Add,
	})
	if err != nil {
		return err
	}

	return o.eventdeletter.DeleteByID(ctx, e.ID)
}
