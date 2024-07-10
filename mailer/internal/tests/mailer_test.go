//go:build integration

package tests

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"strconv"
	"testing"
	"time"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v3/mailer"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/proto"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/app"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/cfg"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/platform/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/rate"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/subscriber"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/transport/nats/subscriber/sub"
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

const (
	testHost = "localhost"
	testPort = "33200"
)

func TestAppGotSubscribersChangedEvent(t *testing.T) {
	const (
		subject    = "subscribers-changed"
		collection = "subscribers"
	)

	type args struct {
		subject string
		event   *sub.SubscriberChangedEvent
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T, db *mongo.Database)
		want    *subscriber.Subscriber
		wantErr bool
	}{
		{
			name: "Should save subscriber to DB",
			args: args{
				subject: subject,
				event: &sub.SubscriberChangedEvent{
					Email: "testinsert@test.com",
				},
			},
			setup: func(*testing.T, *mongo.Database) {},
			want: &subscriber.Subscriber{
				Email: "testinsert@test.com",
			},
			wantErr: false,
		},
		{
			name: "Should not delete subscriber from the DB when it's not present",
			args: args{
				subject: subject,
				event: &sub.SubscriberChangedEvent{
					Email:   "test@test.com",
					Deleted: true,
				},
			},
			setup:   func(*testing.T, *mongo.Database) {},
			wantErr: true,
		},
		{
			name: "Should delete subscriber from the DB",
			args: args{
				subject: subject,
				event: &sub.SubscriberChangedEvent{
					Email:   "test@test.com",
					Deleted: true,
				},
			},
			setup: func(t *testing.T, db *mongo.Database) {
				t.Helper()
				_, err := db.Collection(collection).
					InsertOne(context.Background(), subscriber.Subscriber{
						Email: "test@test.com",
					})
				require.NoError(t, err)
			},
			wantErr: false,
		},
		{
			name: "Should not allow to do request with the empty body",
			args: args{
				subject: subject,
				event:   &sub.SubscriberChangedEvent{},
			},
			setup:   func(*testing.T, *mongo.Database) {},
			wantErr: true,
		},
		{
			name: "Should not allow to do request with the nil body",
			args: args{
				subject: subject,
				event:   nil,
			},
			setup:   func(*testing.T, *mongo.Database) {},
			wantErr: true,
		},
	}

	cfg := mustNewTestConfig(t)
	nc := mustNewNats(t, cfg.NatsURL)
	app := app.New(cfg, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	require.NoError(t, app.Run())
	mongo, err := db.NewConn(context.Background(), cfg.MongoURL)
	require.NoError(t, err)
	db := mongo.GetDB()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		app.Stop()
		nc.Close()
		require.NoError(t, mongo.Close(ctx))
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			t.Cleanup(func() {
				ctx, cancel = context.WithTimeout(context.Background(), time.Second*1)
				defer cancel()
				_, err := db.Collection(collection).DeleteMany(ctx, bson.D{})
				require.NoError(t, err)
			})

			bytes, err := json.Marshal(tt.args.event)
			require.NoError(t, err)

			_, err = nc.RequestWithContext(ctx, subject, bytes)
			require.NoError(t, err)

			subscriber := new(subscriber.Subscriber)

			require.EventuallyWithT(t, func(*assert.CollectT) {
				err = db.Collection(collection).
					FindOne(ctx, bson.D{}).
					Decode(subscriber)

				if tt.wantErr {
					require.Empty(t, subscriber)
					require.Error(t, err)
					return
				}

				got := *subscriber
				if tt.want == nil {
					require.Empty(t, subscriber)
					return
				}

				require.NoError(t, err)
				require.Equal(t, tt.want.Email, got.Email)
			}, time.Second, time.Millisecond*100)
		})
	}
}

func TestAppGotExchangeEvent(t *testing.T) {
	const (
		subject    = "rate-fetched"
		collection = "rate"
	)

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
				subject: subject,
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
				subject: subject,
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
		{
			name: "Should not allow to make request with a nil body",
			args: args{
				subject: subject,
				event:   nil,
			},
			wantErr: true,
		},
		{
			name: "Should not allow to make request with an empty body",
			args: args{
				subject: subject,
				event:   &pb.ExchangeFetchedEvent{},
			},
			wantErr: true,
		},
	}

	cfg := mustNewTestConfig(t)
	nc := mustNewNats(t, cfg.NatsURL)
	app := app.New(cfg, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	require.NoError(t, app.Run())
	mongo, err := db.NewConn(context.Background(), cfg.MongoURL)
	require.NoError(t, err)
	db := mongo.GetDB()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		nc.Close()
		app.Stop()
		require.NoError(t, mongo.Close(ctx))
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			t.Cleanup(func() {
				ctx, cancel = context.WithTimeout(context.Background(), time.Second*1)
				defer cancel()
				_, err := db.Collection(collection).DeleteMany(ctx, bson.D{})
				require.NoError(t, err)
			})

			bytes, err := proto.Marshal(tt.args.event)
			require.NoError(t, err)
			_, err = nc.Request(tt.args.subject, bytes, time.Second)
			require.NoError(t, err)

			exchange := new(rate.Exchange)
			err = db.Collection(collection).
				FindOne(ctx, bson.D{}).
				Decode(exchange)

			if tt.wantErr {
				require.Empty(t, exchange)
				require.Error(t, err)
				return
			}

			got := *exchange
			want := tt.args.event

			require.NoError(t, err)
			require.InDelta(t, want.GetRate(), got.Rate, 2)
			require.Equal(t, want.GetFrom(), got.From)
			require.Equal(t, want.GetTo(), got.To)
		})
	}
}

func mustNewNats(t *testing.T, url string) *nats.Conn {
	t.Helper()
	nc, err := nats.Connect(url)
	require.NoError(t, err, "Failed to connect to NATS")
	return nc
}

func mustNewTestConfig(t *testing.T) cfg.Config {
	t.Helper()
	return cfg.Config{
		MailerFrom:     mustGetEnv(t, mailerSMTPFromEnvKey),
		MailerHost:     mustGetEnv(t, mailerSMTPHostEnvKey),
		MailerPassword: mustGetEnv(t, mailerSMTPPortEnvKey),
		NatsURL:        mustGetEnv(t, mailerTestNatsURL),
		MailerPort:     mustGetIntEnv(t, mailerSMTPPortEnvKey),
		MongoURL:       mustGetEnv(t, mailerTestMongoURL),
		Port:           testPort,
		Host:           testHost,
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
