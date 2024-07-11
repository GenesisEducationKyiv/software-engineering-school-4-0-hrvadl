//go:build integration

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/app"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/cfg"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/event"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/subscriber"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/transport/nats/subscriber/sub"
)

const (
	subTestDSNEnvKey     = "SUB_TEST_DSN"
	subTestNatsURLEnvKey = "SUB_TEST_NATS_URL"
)

const (
	testHost = "localhost"
	testPort = "33200"
)

func TestCompensate(t *testing.T) {
	const subject = "subscribers-changed-failed"

	tests := []struct {
		name  string
		event sub.SubscriberChangedEvent
		setup func(t *testing.T, db *sqlx.DB)
		want  *subscriber.Subscriber
	}{
		{
			name: "Should not compensate deletion with insertion when row is already there",
			event: sub.SubscriberChangedEvent{
				Type:  event.Deleted,
				Email: "mail@mail.com",
			},
			setup: cleanup,
			want: &subscriber.Subscriber{
				Email: "mail@mail.com",
			},
		},
		{
			name: "Should compensate deletion with insertion",
			event: sub.SubscriberChangedEvent{
				Type:  event.Deleted,
				Email: "mail@mail.com",
			},
			setup: cleanup,
			want: &subscriber.Subscriber{
				Email: "mail@mail.com",
			},
		},
		{
			name: "Should not compensate insertion with deletion when no row available",
			event: sub.SubscriberChangedEvent{
				Type:  event.Added,
				Email: "mail@mail.com",
			},
			setup: func(*testing.T, *sqlx.DB) {
			},
			want: nil,
		},
		{
			name: "Should compensate insertion with deletion",
			event: sub.SubscriberChangedEvent{
				Type:  event.Added,
				Email: "mail@mail.com",
			},
			setup: func(t *testing.T, db *sqlx.DB) {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				_, err := db.ExecContext(
					ctx,
					"INSERT INTO subscribers (email) VALUES (?)",
					"mail@mail.com",
				)
				require.NoError(t, err)
			},
			want: nil,
		},
	}

	cfg := cfg.Config{
		Dsn:     mustGetEnv(t, subTestDSNEnvKey),
		Port:    testPort,
		Host:    testHost,
		NatsURL: mustGetEnv(t, subTestNatsURLEnvKey),
	}

	nc := mustNewNats(t, cfg.NatsURL)
	db, err := db.NewConn(cfg.Dsn)
	require.NoError(t, err)

	buf := bytes.NewBufferString("")
	app := app.New(cfg, slog.New(slog.NewTextHandler(buf, nil)))
	go app.MustRun()

	require.EventuallyWithT(t, func(*assert.CollectT) {
		checkPortBusy(t, cfg.Host, cfg.Port)
	}, time.Second, time.Millisecond*100)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanup(t, db)
			})

			tt.setup(t, db)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			bytes, err := json.Marshal(tt.event)
			require.NoError(t, err)

			_, err = nc.RequestWithContext(ctx, subject, bytes)
			require.NoError(t, err)
			var got subscriber.Subscriber
			_ = db.GetContext(
				ctx,
				&got,
				"SELECT * FROM subscribers WHERE email = (?)",
				tt.event.Email,
			)

			if tt.want == nil {
				require.Empty(t, got)
				return
			}

			require.Equal(t, tt.want.Email, got.Email)
		})
	}
}

func cleanup(t *testing.T, db *sqlx.DB) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := db.ExecContext(ctx, "DELETE FROM subscribers")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, "DELETE FROM events")
	require.NoError(t, err)
}

func mustGetEnv(t *testing.T, key string) string {
	t.Helper()
	env := os.Getenv(key)
	require.NotEmpty(t, env)
	return env
}

func mustNewNats(t *testing.T, url string) *nats.Conn {
	t.Helper()
	nc, err := nats.Connect(url)
	require.NoError(t, err, "Failed to connect to NATS")
	return nc
}

func checkPortBusy(t *testing.T, host, port string) {
	t.Helper()
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	require.NoError(t, err)
	require.NotEmpty(t, conn)
}
