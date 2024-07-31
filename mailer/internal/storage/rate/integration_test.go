//go:build integration

package rate

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/platform/db"
)

const mailerMongoTestURLEnvKey = "MAILER_MONGO_TEST_URL"

func TestRepositoryGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	client, err := db.NewConn(ctx, mustGetEnv(t, mailerMongoTestURLEnvKey))
	require.NoError(t, err)
	db := client.GetDB()
	r := &Repository{db: db}
	t.Cleanup(func() {
		require.NoError(t, client.Close(context.Background()))
	})

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T, db *mongo.Database)
		want    *Exchange
		wantErr bool
	}{
		{
			name: "Should return error when not found",
			args: args{
				ctx: context.Background(),
			},
			setup: func(*testing.T, *mongo.Database) {
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should return error when timed out",
			args: args{
				ctx: newImmediateCtx(),
			},
			setup: func(t *testing.T, db *mongo.Database) {
				t.Helper()
				insert(t, db, Exchange{
					ID:   primitive.ObjectID{0, 0, 1},
					From: "USD",
					To:   "UAH",
					Rate: 25.5,
				})
			},
			want: &Exchange{
				ID:   primitive.ObjectID{0, 0, 1},
				From: "USD",
				To:   "UAH",
				Rate: 25.5,
			},
			wantErr: true,
		},
		{
			name: "Should return rate",
			args: args{
				ctx: context.Background(),
			},
			setup: func(t *testing.T, db *mongo.Database) {
				t.Helper()
				insert(t, db, Exchange{
					ID:   primitive.ObjectID{0, 0, 1},
					From: "USD",
					To:   "UAH",
					Rate: 25.5,
				})
			},
			want: &Exchange{
				ID:   primitive.ObjectID{0, 0, 1},
				From: "USD",
				To:   "UAH",
				Rate: 25.5,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				deleteAll(t, db)
			})

			tt.setup(t, db)
			got, err := r.Get(tt.args.ctx)
			if tt.wantErr {
				require.Empty(t, got)
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestRepositoryReplace(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	client, err := db.NewConn(ctx, mustGetEnv(t, mailerMongoTestURLEnvKey))
	require.NoError(t, err)
	db := client.GetDB()
	r := &Repository{db: db}
	t.Cleanup(func() {
		require.NoError(t, client.Close(context.Background()))
	})

	type args struct {
		ctx  context.Context
		rate Exchange
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T, db *mongo.Database)
		wantErr bool
	}{
		{
			name:    "Should return err when timed out",
			wantErr: true,
			setup:   func(*testing.T, *mongo.Database) {},
			args: args{
				ctx: newImmediateCtx(),
				rate: Exchange{
					ID:   primitive.ObjectID{0, 0, 1},
					From: "USD",
					To:   "UAH",
					Rate: 42.2,
				},
			},
		},
		{
			name:    "Should save rate when db is empty",
			wantErr: false,
			setup:   func(*testing.T, *mongo.Database) {},
			args: args{
				ctx: context.Background(),
				rate: Exchange{
					ID:   primitive.ObjectID{0, 0, 1},
					From: "USD",
					To:   "UAH",
					Rate: 42.2,
				},
			},
		},
		{
			name:    "Should replace rate when it's already present",
			wantErr: false,
			setup: func(t *testing.T, db *mongo.Database) {
				t.Helper()
				insert(t, db, Exchange{
					ID:   primitive.ObjectID{0, 0, 1},
					From: "USD",
					To:   "UAH",
					Rate: 42.2,
				})
			},
			args: args{
				ctx: context.Background(),
				rate: Exchange{
					ID:   primitive.ObjectID{0, 0, 1},
					From: "USD",
					To:   "UAH",
					Rate: 42.2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				deleteAll(t, db)
			})

			tt.setup(t, db)
			err := r.Replace(tt.args.ctx, tt.args.rate)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			cur, err := db.Collection("rate").
				Find(tt.args.ctx, bson.D{{Key: "rate", Value: tt.args.rate.Rate}})
			require.NoError(t, err)

			var results []Exchange
			require.NoError(t, cur.All(tt.args.ctx, &results))
			require.Len(t, results, 1)

			rate := results[0]
			require.Equal(t, tt.args.rate, rate)
		})
	}
}

func mustGetEnv(t *testing.T, key string) string {
	t.Helper()
	env := os.Getenv(key)
	require.NotEmpty(t, env)
	return env
}

func deleteAll(t *testing.T, db *mongo.Database) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	_, err := db.Collection("rate").DeleteMany(ctx, bson.D{})
	require.NoError(t, err)
}

func insert(t *testing.T, db *mongo.Database, rate Exchange) {
	t.Helper()
	_, err := db.Collection("rate").InsertOne(context.Background(), rate)
	require.NoError(t, err)
}

func newImmediateCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	return ctx
}
