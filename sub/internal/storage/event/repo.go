package event

import (
	"context"
	"database/sql"
	"fmt"
)

func NewRepo(db DataSource) *Repo {
	return &Repo{
		db: db,
	}
}

type DataSource interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Repo struct {
	db DataSource
}

func (s *Repo) Save(ctx context.Context, event Event) error {
	const query = "INSERT INTO events (type,payload) VALUES (?, ?)"

	if _, err := s.db.ExecContext(ctx, query, event.Type, event.Payload); err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}

	return nil
}
