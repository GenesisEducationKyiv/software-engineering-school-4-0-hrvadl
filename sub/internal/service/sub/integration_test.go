//go:build integration

package sub

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/validator"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/subscriber"
)

const testDSNEnvKey = "SUB_TEST_DSN"

func TestMain(t *testing.M) {
	code := t.Run()
	dsn := os.Getenv(testDSNEnvKey)

	db, err := db.NewConn(dsn)
	if err != nil {
		panic("failed to connect to test db")
	}

	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := db.Exec("DELETE FROM subscribers"); err != nil {
		panic("failed to cleanup")
	}

	os.Exit(code)
}

func TestServiceSubscribe(t *testing.T) {
	type args struct {
		ctx context.Context
		sub subscriber.Subscriber
	}
	testCases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should subscribe correctly",
			args: args{
				ctx: context.Background(),
				sub: subscriber.Subscriber{Email: "test@mail.com"},
			},
			wantErr: false,
		},
		{
			name: "Should not subscribe correctly when it takes too long",
			args: args{
				ctx: newImmediateCtx(),
				sub: subscriber.Subscriber{Email: "test1111@mail.com"},
			},
			wantErr: true,
		},
		{
			name: "Should not subscribe when email is empty",
			args: args{
				ctx: context.Background(),
				sub: subscriber.Subscriber{},
			},
			wantErr: true,
		},
		{
			name: "Should not subscribe when email is incorrect",
			args: args{
				ctx: context.Background(),
				sub: subscriber.Subscriber{Email: "tetmail.com"},
			},
			wantErr: true,
		},
	}

	dsn := os.Getenv(testDSNEnvKey)
	require.NotZero(t, dsn, "test DSN can not be empty")
	dbConn, err := db.NewConn(dsn)
	require.NoError(t, err, "Failed to connect to db")
	t.Cleanup(func() {
		require.NoError(t, dbConn.Close(), "Failed to close DB")
	})

	rs := subscriber.NewRepo(db.NewWithTx(dbConn))
	v := validator.NewStdlib()
	s := NewService(rs, v)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			id, err := s.Subscribe(tt.args.ctx, tt.args.sub)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			t.Cleanup(func() {
				cleanupSub(t, dbConn, id)
			})

			require.NoError(t, err)
			require.NotZero(t, id)
		})
	}
}

func TestServiceSubscribeTwice(t *testing.T) {
	type args struct {
		ctx context.Context
		sub subscriber.Subscriber
	}
	testCases := []struct {
		name string
		args args
	}{
		{
			name: "Should not subscribe twice",
			args: args{
				ctx: context.Background(),
				sub: subscriber.Subscriber{Email: "test@mail.com"},
			},
		},
		{
			name: "Should not subscribe twice",
			args: args{
				ctx: context.Background(),
				sub: subscriber.Subscriber{Email: "testnew@mail.com"},
			},
		},
	}

	dsn := os.Getenv(testDSNEnvKey)
	require.NotZero(t, dsn, "test DSN can not be empty")
	dbConn, err := db.NewConn(dsn)
	require.NoError(t, err, "Failed to connect to db")
	t.Cleanup(func() {
		require.NoError(t, dbConn.Close(), "Failed to close DB")
	})

	rs := subscriber.NewRepo(db.NewWithTx(dbConn))
	v := validator.NewStdlib()
	s := NewService(rs, v)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			id, err := s.Subscribe(tt.args.ctx, tt.args.sub)
			t.Cleanup(func() {
				cleanupSub(t, dbConn, id)
			})

			require.NoError(t, err)
			require.NotZero(t, id)

			id, err = s.Subscribe(tt.args.ctx, tt.args.sub)
			require.Error(t, err)
			require.Zero(t, id)
		})
	}
}

func newImmediateCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	return ctx
}

func cleanupSub(t *testing.T, db *sqlx.DB, id int64) {
	t.Helper()
	_, err := db.Exec("DELETE FROM subscribers WHERE id = ?", id)
	require.NoError(t, err, "Failed to clean up subscriber")
}
