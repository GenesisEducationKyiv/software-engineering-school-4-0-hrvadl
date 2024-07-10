package event

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewGetter(db *sqlx.DB) *Getter {
	return &Getter{
		db: db,
	}
}

type Getter struct {
	db *sqlx.DB
}

func (g *Getter) GetAll(ctx context.Context) ([]Event, error) {
	const query = "SELECT * FROM events"
	var events []Event

	if err := g.db.SelectContext(ctx, &events, query); err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	return events, nil
}
