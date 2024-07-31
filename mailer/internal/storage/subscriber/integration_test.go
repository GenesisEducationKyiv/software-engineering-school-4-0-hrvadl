//go:build integration

package subscriber

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

func TestRepositoryGetAll(t *testing.T) {
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
		want    []Subscriber
		wantErr bool
	}{
		{
			name: "Should return one subscriber",
			args: args{
				ctx: context.Background(),
			},
			setup: func(t *testing.T, db *mongo.Database) {
				t.Helper()
				insert(t, db, Subscriber{
					ID:    primitive.ObjectID{0, 0, 1},
					Email: "test@test.com",
				})
			},
			want: []Subscriber{
				{
					ID:    primitive.ObjectID{0, 0, 1},
					Email: "test@test.com",
				},
			},
		},
		{
			name: "Should return one multiple subscriber",
			args: args{
				ctx: context.Background(),
			},
			setup: func(t *testing.T, db *mongo.Database) {
				t.Helper()
				insert(t, db, Subscriber{
					ID:    primitive.ObjectID{0, 0, 1},
					Email: "test@test.com",
				})
				insert(t, db, Subscriber{
					ID:    primitive.ObjectID{0, 1, 1},
					Email: "test2@test.com",
				})
			},
			want: []Subscriber{
				{
					ID:    primitive.ObjectID{0, 0, 1},
					Email: "test@test.com",
				},
				{
					ID:    primitive.ObjectID{0, 1, 1},
					Email: "test2@test.com",
				},
			},
		},
		{
			name: "Should return error when timed out",
			args: args{
				ctx: newImmediateCtx(),
			},
			setup: func(t *testing.T, db *mongo.Database) {
				t.Helper()
				insert(t, db, Subscriber{
					ID:    primitive.ObjectID{0, 0, 1},
					Email: "test@test.com",
				})
				insert(t, db, Subscriber{
					ID:    primitive.ObjectID{0, 1, 1},
					Email: "test2@test.com",
				})
			},
			wantErr: true,
		},
		{
			name: "Should not return err when db is empty",
			args: args{
				ctx: context.Background(),
			},
			setup: func(*testing.T, *mongo.Database) {
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
			got, err := r.GetAll(tt.args.ctx)
			if tt.wantErr {
				t.Log(len(got))
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestRepositorySave(t *testing.T) {
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
		sub Subscriber
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should save to repository correctly",
			args: args{
				sub: Subscriber{ID: primitive.ObjectID{0, 0, 1}, Email: "test@test.com"},
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Should save to repository correctly",
			args: args{
				sub: Subscriber{ID: primitive.ObjectID{0, 0, 1, 1}, Email: "test2@test.com"},
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Should return error when timed out",
			args: args{
				sub: Subscriber{ID: primitive.ObjectID{0, 0, 1, 1}, Email: "test2@test.com"},
				ctx: newImmediateCtx(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				deleteAll(t, db)
			})

			err := r.Save(tt.args.ctx, tt.args.sub)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			cur, err := db.Collection("subscribers").Find(tt.args.ctx, bson.D{{}})
			require.NoError(t, err)

			var sub []Subscriber
			err = cur.All(context.Background(), &sub)
			require.NoError(t, err)

			require.Len(t, sub, 1)
			require.Equal(t, tt.args.sub, sub[0])
		})
	}
}

func TestRepositoryDelete(t *testing.T) {
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
		sub Subscriber
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T, db *mongo.Database)
		wantErr bool
	}{
		{
			name: "Should return error if subscriber does not exist",
			args: args{
				ctx: context.Background(),
				sub: Subscriber{Email: "dont@exist.com"},
			},
			setup:   func(*testing.T, *mongo.Database) {},
			wantErr: true,
		},
		{
			name: "Should return error when request timed out",
			args: args{
				ctx: newImmediateCtx(),
				sub: Subscriber{Email: "test1@test.com"},
			},
			setup: func(t *testing.T, db *mongo.Database) {
				t.Helper()
				_, err := db.Collection("subscribers").InsertOne(context.Background(), Subscriber{
					Email: "test1@test.com",
				})
				require.NoError(t, err)
			},
			wantErr: true,
		},
		{
			name: "Should delete subscriber",
			args: args{
				ctx: context.Background(),
				sub: Subscriber{Email: "test@test.com"},
			},
			setup: func(t *testing.T, db *mongo.Database) {
				t.Helper()
				_, err := db.Collection("subscribers").InsertOne(context.Background(), Subscriber{
					Email: "test@test.com",
				})
				require.NoError(t, err)
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
			err := r.Delete(tt.args.ctx, tt.args.sub)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
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
	_, err := db.Collection("subscribers").DeleteMany(ctx, bson.D{})
	require.NoError(t, err)
}

func insert(t *testing.T, db *mongo.Database, sub Subscriber) {
	t.Helper()
	_, err := db.Collection("subscribers").InsertOne(context.Background(), sub)
	require.NoError(t, err)
}

func newImmediateCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	return ctx
}
