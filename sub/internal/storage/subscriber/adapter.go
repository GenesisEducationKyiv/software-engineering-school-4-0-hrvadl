package subscriber

import "context"

func NewCompensateAdapter(r *Repo) *CompensateAdapter {
	return &CompensateAdapter{
		repo: r,
	}
}

type CompensateAdapter struct {
	repo *Repo
}

func (c *CompensateAdapter) Save(ctx context.Context, sub Subscriber) (int64, error) {
	return c.repo.CompensateSave(ctx, sub)
}

func (c *CompensateAdapter) DeleteByEmail(ctx context.Context, email string) error {
	return c.repo.CompensateDeleteByEmail(ctx, email)
}

func (c *CompensateAdapter) GetByEmail(ctx context.Context, email string) (*Subscriber, error) {
	return c.repo.GetByEmail(ctx, email)
}
