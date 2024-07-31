package subscriber

import (
	"context"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/transaction"
)

func NewCompensateAdapter(r *Repo, tx *transaction.Manager) *CompensateAdapter {
	return &CompensateAdapter{
		repo: r,
		tx:   tx,
	}
}

type CompensateAdapter struct {
	repo *Repo
	tx   *transaction.Manager
}

func (c *CompensateAdapter) Save(ctx context.Context, sub Subscriber) (int64, error) {
	var id int64
	err := c.tx.WithTx(ctx, func(ctx context.Context) error {
		var err error
		id, err = c.repo.Save(ctx, sub)
		return err
	})
	return id, err
}

func (c *CompensateAdapter) DeleteByEmail(ctx context.Context, email string) error {
	return c.tx.WithTx(ctx, func(ctx context.Context) error {
		return c.repo.DeleteByEmail(ctx, email)
	})
}

func (c *CompensateAdapter) GetByEmail(ctx context.Context, email string) (*Subscriber, error) {
	return c.repo.GetByEmail(ctx, email)
}
