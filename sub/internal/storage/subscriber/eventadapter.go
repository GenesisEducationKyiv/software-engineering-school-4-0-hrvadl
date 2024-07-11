package subscriber

import (
	"context"
	"errors"

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
		var err error
		id, err = c.repo.Save(ctx, sub)
		evErr := c.events.Save(ctx, event.Event{Payload: sub.Email, Type: event.Added})
		return errors.Join(err, evErr)
	})
	return id, err
}

func (c *WithEventAdapter) DeleteByEmail(ctx context.Context, email string) error {
	return c.tx.WithTx(ctx, func(ctx context.Context) error {
		err := c.repo.DeleteByEmail(ctx, email)
		evErr := c.events.Save(ctx, event.Event{Payload: email, Type: event.Deleted})
		return errors.Join(err, evErr)
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
