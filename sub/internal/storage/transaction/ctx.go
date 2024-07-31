package transaction

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type key string

const contextKey = key("string")

func AddToContext(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, contextKey, tx)
}

func FromContext(ctx context.Context) (*sqlx.Tx, error) {
	tx, ok := ctx.Value(contextKey).(*sqlx.Tx)
	if !ok {
		return nil, ErrFailedGetTx
	}

	return tx, nil
}
