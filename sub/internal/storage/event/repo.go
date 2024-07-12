package event

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
)

func NewRepo(db *db.Tx) *Repo {
	return &Repo{
		db: db,
	}
}

type Repo struct {
	db *db.Tx
}

func (s *Repo) Save(ctx context.Context, event Event) error {
	const query = "INSERT INTO events (type,payload) VALUES (?, ?)"

	if _, err := s.db.ExecContext(ctx, query, event.Type, event.Payload); err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}

	return nil
}
