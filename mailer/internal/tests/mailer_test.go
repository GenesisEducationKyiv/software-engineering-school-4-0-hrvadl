//go:build integration

package tests

import (
	"bytes"
	"context"
	"log/slog"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v3/mailer"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/proto"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/app"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/cfg"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/platform/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/rate"
)

const (
	mailerSMTPHostEnvKey     = "MAILER_TEST_SMTP_HOST"
	mailerSMTPPasswordEnvKey = "MAILER_TEST_SMTP_PASSWORD" // #nosec G101
	mailerSMTPFromEnvKey     = "MAILER_TEST_SMTP_FROM"
	mailerSMTPPortEnvKey     = "MAILER_TEST_SMTP_PORT"
	mailerTestAPIPortEnvKey  = "MAILER_TEST_API_PORT"
	mailerTestNatsURL        = "MAILER_NATS_TEST_URL"
	mailerTestMongoURL       = "MAILER_MONGO_TEST_URL"
)

// NOTE: I'm not doing component test using Resend mailer.
// It has quota and I don't wanna exceed it.
func TestAppGotExchangeEvent(t *testing.T) {
	const topic = "rate-fetched"
	type args struct {
		subject string
		event   *pb.ExchangeFetchedEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should save exchange to DB",
			args: args{
				subject: topic,
				event: &pb.ExchangeFetchedEvent{
					EventID:   "1",
					EventType: "rate-fetched",
					From:      "USD",
					To:        "UAH",
					Rate:      32.2,
				},
			},
		},
		{
			name: "Should ignore unknown event type",
			args: args{
				subject: "rate-fetched",
				event: &pb.ExchangeFetchedEvent{
					EventID:   "1",
					EventType: "unknown",
					From:      "USD",
					To:        "UAH",
					Rate:      32.3,
				},
			},
			wantErr: true,
		},
	}

	const (
		testHost = "localhost"
		testPort = "33200"
	)

	cfg := cfg.Config{
		MailerFrom:     mustGetEnv(t, mailerSMTPFromEnvKey),
		MailerHost:     mustGetEnv(t, mailerSMTPHostEnvKey),
		MailerPassword: mustGetEnv(t, mailerSMTPPortEnvKey),
		NatsURL:        mustGetEnv(t, mailerTestNatsURL),
		MailerPort:     mustGetIntEnv(t, mailerSMTPPortEnvKey),
		MongoURL:       mustGetEnv(t, mailerTestMongoURL),
		Port:           testPort,
		Host:           testHost,
	}

	nc, err := nats.Connect(cfg.NatsURL)
	require.NoError(t, err, "Failed to connect to NATS")
	js, err := jetstream.New(nc)
	require.NoError(t, err, "Failed to connect to JetStream")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err = js.CreateStream(
		ctx,
		jetstream.StreamConfig{
			Name:     "DebeziumStream",
			Subjects: []string{"subscribers-changed"},
		},
	)
	require.NoError(t, err, "Failed to create JetStream")

	buf := bytes.NewBuffer([]byte{})
	app := app.New(cfg, slog.New(slog.NewTextHandler(buf, nil)))
	require.NoError(t, app.Run())
	mongo, err := db.NewConn(context.Background(), cfg.MongoURL)
	require.NoError(t, err)
	db := mongo.GetDB()

	t.Cleanup(func() {
		app.Stop()
		require.NoError(t, mongo.Close(context.Background()))
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				_, err := db.Collection("rate").DeleteMany(context.Background(), bson.D{})
				require.NoError(t, err)
			})

			bytes, err := proto.Marshal(tt.args.event)
			require.NoError(t, err)
			_, err = nc.Request(tt.args.subject, bytes, time.Second)
			require.NoError(t, err)

			require.EventuallyWithT(t, func(*assert.CollectT) {
				ex := new(rate.Exchange)
				err := db.Collection("rate").
					FindOne(context.Background(), bson.D{}).
					Decode(ex)

				if tt.wantErr {
					require.Empty(t, ex)
					require.Error(t, err)
					return
				}

				got := *ex
				want := tt.args.event

				require.NoError(t, err)
				require.InDelta(t, want.GetRate(), got.Rate, 2)
				require.Equal(t, want.GetFrom(), got.From)
				require.Equal(t, want.GetTo(), got.To)
			}, time.Second, 100*time.Millisecond)
		})
	}
}

func mustGetIntEnv(t *testing.T, key string) int {
	t.Helper()
	i, err := strconv.Atoi(mustGetEnv(t, key))
	require.NoError(t, err)
	return i
}

func mustGetEnv(t *testing.T, key string) string {
	t.Helper()
	env := os.Getenv(key)
	require.NotEmpty(t, env)
	return env
}

func checkPortBusy(t *testing.T, host string, port string) {
	t.Helper()
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	require.NoError(t, err)
	require.NotEmpty(t, conn)
	require.NoError(t, conn.Close())
}
