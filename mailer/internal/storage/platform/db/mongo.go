package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	operation = "mongo connect"
	db        = "notifier"
)

func NewConn(ctx context.Context, url string) (*Conn, error) {
	bsonOpts := &options.BSONOptions{
		UseJSONStructTags: true,
		NilMapAsEmpty:     true,
		NilSliceAsEmpty:   true,
	}

	clientOpts := options.Client().
		ApplyURI(url).
		SetBSONOptions(bsonOpts)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", operation, ErrFailedConnect, err)
	}
	return &Conn{
		client: client,
		db:     client.Database(db),
	}, nil
}

type Conn struct {
	client *mongo.Client
	db     *mongo.Database
}

func (c *Conn) GetDB() *mongo.Database {
	return c.db
}

func (c *Conn) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}
