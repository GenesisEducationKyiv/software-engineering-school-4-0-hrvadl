package rate

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collection = "rate"
	operation  = "rate repo"
)

func NewRepo(db *mongo.Database) *Repository {
	return &Repository{
		db: db,
	}
}

type Repository struct {
	db *mongo.Database
}

func (r *Repository) Get(ctx context.Context) (*Exchange, error) {
	var ex Exchange
	if err := r.db.Collection(collection).FindOne(ctx, bson.D{}).Decode(&ex); err != nil {
		return nil, fmt.Errorf("%s: failed to get rate: %w", operation, err)
	}

	return &ex, nil
}

func (r *Repository) Replace(ctx context.Context, rate Exchange) error {
	doc, err := r.Get(ctx)
	switch {
	case err == nil:
		return r.replace(ctx, *doc, rate)
	case errors.Is(err, mongo.ErrNoDocuments):
		return r.save(ctx, rate)
	default:
		return fmt.Errorf("%s: failed to replace rate %w", operation, err)
	}
}

func (r *Repository) replace(ctx context.Context, old Exchange, replace Exchange) error {
	if res := r.db.Collection(collection).FindOneAndReplace(ctx, old, replace); res.Err() != nil {
		return fmt.Errorf("%s: %w", operation, res.Err())
	}
	return nil
}

func (r *Repository) save(ctx context.Context, rate Exchange) error {
	if _, err := r.db.Collection(collection).InsertOne(ctx, rate); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}
