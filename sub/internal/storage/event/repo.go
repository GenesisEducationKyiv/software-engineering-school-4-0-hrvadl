package event

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewRepo(tx *sqlx.Tx) *Repo {
	return &Repo{
		tx: tx,
	}
}

type Repo struct {
	tx *sqlx.Tx
}

func (s *Repo) Save(ctx context.Context, event Event) error {
	const query = "INSERT INTO events (type,payload) VALUES (?, ?)"
	if _, err := s.tx.ExecContext(ctx, query, event.Type, event.Payload); err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}

	return nil
}
