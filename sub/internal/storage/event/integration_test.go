//go:build integration

package event

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
)

const dsnEnvKey = "SUB_TEST_DSN"

func TestRepoSave(t *testing.T) {
	type args struct {
		ctx   context.Context
		event Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Should save delete event correctly",
			wantErr: false,
			args: args{
				ctx: context.Background(),
				event: Event{
					Type:    Deleted,
					Payload: "test@test.com",
				},
			},
		},
		{
			name:    "Should save insert event correctly",
			wantErr: false,
			args: args{
				ctx: context.Background(),
				event: Event{
					Type:    Added,
					Payload: "test@test.com",
				},
			},
		},
		{
			name:    "Should return error when timed out",
			wantErr: true,
			args: args{
				ctx: newImmediateCtx(),
				event: Event{
					Type:    Added,
					Payload: "test@test.com",
				},
			},
		},
	}

	dbConn, err := db.NewConn(mustGetEnv(t, dsnEnvKey))
	require.NoError(t, err)
	txDB := db.NewWithTx(dbConn)
	repo := NewRepo(txDB)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanup(t, txDB)
			})

			err := repo.Save(tt.args.ctx, tt.args.event)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()

			var event Event
			err = txDB.GetContext(
				ctx,
				&event,
				"SELECT * FROM events WHERE payload = (?)",
				tt.args.event.Payload,
			)
			require.NoError(t, err)
		})
	}
}

func newImmediateCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	return ctx
}

func cleanup(t *testing.T, db *db.Tx) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	_, err := db.ExecContext(ctx, "DELETE FROM events")
	require.NoError(t, err)
}

func mustGetEnv(t *testing.T, key string) string {
	t.Helper()
	env := os.Getenv(key)
	require.NotEmpty(t, env)
	return env
}
