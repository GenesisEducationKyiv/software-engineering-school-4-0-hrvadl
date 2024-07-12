package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/transaction"
)

func NewWithTx(db *sqlx.DB) *Tx {
	return &Tx{
		db: db,
	}
}

type Tx struct {
	db *sqlx.DB
}

func (d *Tx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	tx, err := transaction.FromContext(ctx)
	if err != nil {
		return d.db.ExecContext(ctx, query, args...)
	}
	return tx.ExecContext(ctx, query, args...)
}

func (d *Tx) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.GetContext(ctx, dest, query, args...)
}

func (d *Tx) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	return d.db.SelectContext(ctx, dest, query, args...)
}
