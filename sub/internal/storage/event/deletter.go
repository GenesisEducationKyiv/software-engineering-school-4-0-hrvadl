package event

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewDeletter(db *sqlx.DB) *Deletter {
	return &Deletter{
		db: db,
	}
}

type Deletter struct {
	db *sqlx.DB
}

func (g *Deletter) DeleteByID(ctx context.Context, id int) error {
	const query = "DELETE FROM events WHERE id = (?)"

	if _, err := g.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("failed to get events: %w", err)
	}

	return nil
}
