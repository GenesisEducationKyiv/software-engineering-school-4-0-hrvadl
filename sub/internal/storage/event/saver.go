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

func (s *Repo) GetByPayload(ctx context.Context, etype, payload string) (*Event, error) {
	const query = "SELECT * FROM events WHERE type = (?) AND payload = (?)"
	var event Event
	if err := s.tx.GetContext(ctx, &event, query, etype, payload); err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &event, nil
}

func (s *Repo) DeleteByID(ctx context.Context, id int) error {
	const query = "DELETE FROM events WHERE id = (?)"
	if _, err := s.tx.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}
