//go:build !integration

package transaction

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestAddToContext(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		tx  *sqlx.Tx
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "Should add to context correctly",
			args: args{
				ctx: context.Background(),
				tx:  &sqlx.Tx{},
			},
			want: context.WithValue(context.Background(), contextKey, &sqlx.Tx{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := AddToContext(tt.args.ctx, tt.args.tx)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFromContext(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *sqlx.Tx
		wantErr bool
	}{
		{
			name: "Should extract from context correctly",
			args: args{
				ctx: context.WithValue(context.Background(), contextKey, &sqlx.Tx{}),
			},
			want: &sqlx.Tx{},
		},
		{
			name: "Should not extract from context correctly",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := FromContext(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
