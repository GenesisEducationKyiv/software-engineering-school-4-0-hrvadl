package transaction

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

func NewManager(db *sqlx.DB) *Manager {
	return &Manager{
		db: db,
	}
}

type Manager struct {
	db *sqlx.DB
}

func (m *Manager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Join(ErrFailedBeginTx, err)
	}

	ctx = AddToContext(ctx, tx)
	if err := fn(ctx); err != nil {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}
