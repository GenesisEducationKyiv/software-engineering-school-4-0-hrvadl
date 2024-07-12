//go:build integration

package subscriber

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
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

func TestSave(t *testing.T) {
	type args struct {
		ctx context.Context
		sub Subscriber
	}
	testCases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should save subscriber correctly",
			args: args{
				ctx: context.Background(),
				sub: Subscriber{Email: "test1@mail.com"},
			},
			wantErr: false,
		},
		{
			name: "Should not get subscribers correctly when it takes too long",
			args: args{
				ctx: newImmediateCtx(),
			},
			wantErr: true,
		},
	}

	dsn := os.Getenv(testDSNEnvKey)
	require.NotZero(t, dsn, "test DSN can not be empty")
	dbConn, err := db.NewConn(dsn)
	require.NoError(t, err, "Failed to connect to test DB")
	t.Cleanup(func() {
		require.NoError(t, dbConn.Close(), "Failed to close DB")
	})
	r := NewRepo(db.NewWithTx(dbConn))

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			id, err := r.Save(tt.args.ctx, tt.args.sub)
			t.Cleanup(func() {
				cleanupSub(t, dbConn, id)
			})

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotZero(t, id)
		})
	}
}

func TestSaveSubscriberTwice(t *testing.T) {
	type args struct {
		ctx context.Context
		sub Subscriber
	}
	testCases := []struct {
		name string
		args args
	}{
		{
			name: "Should not save subscriber twice",
			args: args{
				ctx: context.Background(),
				sub: Subscriber{Email: "test1@mail.com"},
			},
		},
		{
			name: "Should not save subscriber twice",
			args: args{
				ctx: context.Background(),
				sub: Subscriber{Email: "test@mail.com"},
			},
		},
	}

	dsn := os.Getenv(testDSNEnvKey)
	require.NotZero(t, dsn, "test DSN can not be empty")
	dbConn, err := db.NewConn(dsn)
	require.NoError(t, err, "Failed to connect to test DB")
	t.Cleanup(func() {
		require.NoError(t, dbConn.Close(), "Failed to close DB")
	})
	r := NewRepo(db.NewWithTx(dbConn))

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			id, err := r.Save(tt.args.ctx, tt.args.sub)
			t.Cleanup(func() {
				cleanupSub(t, dbConn, id)
			})

			require.NoError(t, err)
			require.NotZero(t, id)

			id, err = r.Save(tt.args.ctx, tt.args.sub)
			require.Error(t, err)
			require.Zero(t, id)
		})
	}
}

func TestGetSubscribers(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	testCases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should get subscribers correctly",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Should not get subscribers correctly when it takes too long",
			args: args{
				ctx: newImmediateCtx(),
			},
			wantErr: true,
		},
	}

	dsn := os.Getenv(testDSNEnvKey)
	require.NotZero(t, dsn, "test DSN can not be empty")
	dbConn, err := db.NewConn(dsn)
	require.NoError(t, err, "Failed to connect to test DB")
	t.Cleanup(func() {
		require.NoError(t, dbConn.Close(), "Failed to close DB")
	})
	r := NewRepo(db.NewWithTx(dbConn))

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			want := seed(t, r, 30)
			t.Cleanup(func() {
				for _, s := range want {
					cleanupSub(t, dbConn, s.ID)
				}
			})

			got, err := r.Get(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Subset(t, mapSubsToMails(got), mapSubsToMails(want))
		})
	}
}

func seed(t *testing.T, repo *Repo, amount int) []Subscriber {
	t.Helper()

	subs := make([]Subscriber, 0, amount)
	for range amount {
		mail := fmt.Sprintf("mail%v@mail.com", time.Now().Nanosecond())
		sub := Subscriber{Email: mail}
		subs = append(subs, sub)
		id, err := repo.Save(context.Background(), sub)
		sub.ID = id
		require.NoError(t, err)
		require.NotZero(t, id)
	}

	return subs
}

func newImmediateCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	return ctx
}

func mapSubsToMails(s []Subscriber) []string {
	mails := make([]string, 0, len(s))
	for i := range s {
		mails = append(mails, s[i].Email)
	}
	return mails
}

func cleanupSub(t *testing.T, db *sqlx.DB, id int64) {
	t.Helper()
	_, err := db.Exec("DELETE FROM subscribers WHERE id = ?", id)
	require.NoError(t, err, "Failed to clean up subscriber")
}
