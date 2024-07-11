package event

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/transaction"
)

func NewRepo() *Repo {
	return &Repo{}
}

type Repo struct{}

func (s *Repo) Save(ctx context.Context, event Event) error {
	const query = "INSERT INTO events (type,payload) VALUES (?, ?)"

	tx, err := transaction.FromContext(ctx)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, query, event.Type, event.Payload); err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}

	return nil
}
