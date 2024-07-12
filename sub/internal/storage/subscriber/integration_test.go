//go:build integration

package subscriber

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/event"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/transaction"
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

	dsn := mustGetEnv(t, testDSNEnvKey)
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

	dsn := mustGetEnv(t, testDSNEnvKey)
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

	dsn := mustGetEnv(t, testDSNEnvKey)
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

func TestWithEventAdapterSave(t *testing.T) {
	type args struct {
		ctx context.Context
		sub Subscriber
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should save subscriber and create event",
			args: args{
				ctx: context.Background(),
				sub: Subscriber{
					Email: "test@test1.com",
				},
			},
			wantErr: false,
		},
		{
			name: "Should not save subscriber and create event when timed out",
			args: args{
				ctx: newImmediateCtx(),
				sub: Subscriber{
					Email: "test@test1.com",
				},
			},
			wantErr: true,
		},
	}

	dbConn, err := db.NewConn(mustGetEnv(t, testDSNEnvKey))
	require.NoError(t, err)
	txDB := db.NewWithTx(dbConn)
	repo := NewRepo(txDB)
	events := event.NewRepo(txDB)
	tx := transaction.NewManager(dbConn)

	c := &WithEventAdapter{
		repo:   repo,
		events: events,
		tx:     tx,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanup(t, txDB)
			})

			_, err := c.Save(tt.args.ctx, tt.args.sub)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			var sub Subscriber
			err = txDB.GetContext(
				ctx,
				&sub,
				"SELECT * FROM subscribers WHERE email = (?)",
				tt.args.sub.Email,
			)
			require.NoError(t, err)
			require.Equal(t, sub.Email, tt.args.sub.Email)

			var event event.Event
			err = txDB.GetContext(
				ctx,
				&event,
				"SELECT * FROM events WHERE payload = (?) AND type = 'subscriber-added'",
				tt.args.sub.Email,
			)
			require.NoError(t, err)
			require.Equal(t, event.Payload, tt.args.sub.Email)
		})
	}
}

func TestWithEventAdapterDeleteByEmail(t *testing.T) {
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T, db *db.Tx)
		wantErr bool
	}{
		{
			name: "Shoould delete subscriber by email",
			args: args{
				ctx:   context.Background(),
				email: "test@test.com",
			},
			setup: func(t *testing.T, db *db.Tx) {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				_, err := db.ExecContext(
					ctx,
					"INSERT INTO subscribers (email) VALUES (?)",
					"test@test.com",
				)
				require.NoError(t, err)
			},
			wantErr: false,
		},
	}

	dbConn, err := db.NewConn(mustGetEnv(t, testDSNEnvKey))
	require.NoError(t, err)
	txDB := db.NewWithTx(dbConn)
	repo := NewRepo(txDB)
	events := event.NewRepo(txDB)
	tx := transaction.NewManager(dbConn)

	c := &WithEventAdapter{
		repo:   repo,
		events: events,
		tx:     tx,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanup(t, txDB)
			})
			tt.setup(t, txDB)
			err := c.DeleteByEmail(tt.args.ctx, tt.args.email)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			var sub Subscriber
			err = txDB.GetContext(
				ctx,
				&sub,
				"SELECT * FROM subscribers WHERE email = (?)",
				tt.args.email,
			)
			require.Error(t, err)
			require.ErrorIs(t, err, sql.ErrNoRows)

			var event event.Event
			err = txDB.GetContext(
				ctx,
				&event,
				"SELECT * FROM events WHERE payload = (?) AND type = 'subscriber-deleted'",
				tt.args.email,
			)
			require.NoError(t, err)
			require.Equal(t, tt.args.email, event.Payload)
		})
	}
}

func TestWithEventAdapterGetByEmail(t *testing.T) {
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T, db *db.Tx)
		want    *Subscriber
		wantErr bool
	}{
		{
			name: "Should get subscriber correctly",
			args: args{
				ctx:   context.Background(),
				email: "test@test.com",
			},
			setup: func(t *testing.T, db *db.Tx) {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				_, err := db.ExecContext(
					ctx,
					"INSERT INTO subscribers (email) VALUES (?)",
					"test@test.com",
				)
				require.NoError(t, err)
			},
			want: &Subscriber{
				Email: "test@test.com",
			},
		},
		{
			name: "Should return error when subscriber does not exist",
			args: args{
				ctx:   context.Background(),
				email: "test@test.com",
			},
			setup: func(*testing.T, *db.Tx) {
			},
			wantErr: true,
		},
	}

	dbConn, err := db.NewConn(mustGetEnv(t, testDSNEnvKey))
	require.NoError(t, err)
	txDB := db.NewWithTx(dbConn)
	repo := NewRepo(txDB)
	events := event.NewRepo(txDB)
	tx := transaction.NewManager(dbConn)

	c := &WithEventAdapter{
		repo:   repo,
		events: events,
		tx:     tx,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanup(t, txDB)
			})
			tt.setup(t, txDB)
			got, err := c.GetByEmail(tt.args.ctx, tt.args.email)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.Email, got.Email)
		})
	}
}

func TestWithEventAdapterGet(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T, db *db.Tx)
		want    []Subscriber
		wantErr bool
	}{
		{
			name: "Should get all subscribers",
			args: args{
				ctx: context.Background(),
			},
			setup: func(t *testing.T, db *db.Tx) {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				_, err := db.NamedExecContext(
					ctx,
					"INSERT INTO subscribers (email) VALUES (:email)",
					[]Subscriber{
						{Email: "test@test.com"},
					},
				)
				require.NoError(t, err)
			},
			want: []Subscriber{
				{Email: "test@test.com"},
			},
		},
		{
			name: "Should get all subscribers",
			args: args{
				ctx: context.Background(),
			},
			setup: func(t *testing.T, db *db.Tx) {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				_, err := db.NamedExecContext(
					ctx,
					"INSERT INTO subscribers (email) VALUES (:email)",
					[]Subscriber{
						{Email: "test@test.com"},
						{Email: "test1@test.com"},
					},
				)
				require.NoError(t, err)
			},
			want: []Subscriber{
				{Email: "test@test.com"},
				{Email: "test1@test.com"},
			},
		},
	}

	dbConn, err := db.NewConn(mustGetEnv(t, testDSNEnvKey))
	require.NoError(t, err)
	txDB := db.NewWithTx(dbConn)
	repo := NewRepo(txDB)
	events := event.NewRepo(txDB)
	tx := transaction.NewManager(dbConn)

	c := &WithEventAdapter{
		repo:   repo,
		events: events,
		tx:     tx,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t, txDB)
			t.Cleanup(func() {
				cleanup(t, txDB)
			})

			got, err := c.Get(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, mapSubsToMails(tt.want), mapSubsToMails(got))
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

func cleanup(t *testing.T, db *db.Tx) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := db.ExecContext(ctx, "DELETE FROM subscribers")
	require.NoError(t, err, "Failed to clean up subscriber")
	_, err = db.ExecContext(ctx, "DELETE FROM events")
	require.NoError(t, err, "Failed to clean up subscriber")
}

func mustGetEnv(t *testing.T, key string) string {
	t.Helper()
	env := os.Getenv(key)
	require.NotEmpty(t, env)
	return env
}
