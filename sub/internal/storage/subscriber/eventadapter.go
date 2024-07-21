package subscriber

import (
	"context"
	"encoding/json"
	"errors"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/sub"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/event"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/transaction"
)

func NewWithEventAdapter(r *Repo, er *event.Repo, tx *transaction.Manager) *WithEventAdapter {
	return &WithEventAdapter{
		repo:   r,
		events: er,
		tx:     tx,
	}
}

type WithEventAdapter struct {
	repo   *Repo
	events *event.Repo
	tx     *transaction.Manager
}

func (c *WithEventAdapter) Save(ctx context.Context, sub Subscriber) (int64, error) {
	var id int64
	err := c.tx.WithTx(ctx, func(ctx context.Context) error {
		subAdded := &pb.SubscriptionAddedEvent{Email: sub.Email}
		payload, err := json.Marshal(subAdded)
		if err != nil {
			return err
		}

		evErr := c.events.Save(ctx, event.Event{Payload: payload, Type: event.Added})
		id, err = c.repo.Save(ctx, sub)
		return errors.Join(err, evErr)
	})
	return id, err
}

func (c *WithEventAdapter) DeleteByEmail(ctx context.Context, email string) error {
	return c.tx.WithTx(ctx, func(ctx context.Context) error {
		subDeleted := &pb.SubscriptionAddedEvent{Email: email}
		payload, err := json.Marshal(subDeleted)
		if err != nil {
			return err
		}

		evErr := c.events.Save(ctx, event.Event{Payload: payload, Type: event.Deleted})
		return errors.Join(c.repo.DeleteByEmail(ctx, email), evErr)
	})
}

func (c *WithEventAdapter) GetByEmail(ctx context.Context, email string) (*Subscriber, error) {
	var sub *Subscriber
	err := c.tx.WithTx(ctx, func(ctx context.Context) error {
		var err error
		sub, err = c.repo.GetByEmail(ctx, email)
		return err
	})
	return sub, err
}

func (c *WithEventAdapter) Get(ctx context.Context) ([]Subscriber, error) {
	var sub []Subscriber
	err := c.tx.WithTx(ctx, func(ctx context.Context) error {
		var err error
		sub, err = c.repo.Get(ctx)
		return err
	})
	return sub, err
}
