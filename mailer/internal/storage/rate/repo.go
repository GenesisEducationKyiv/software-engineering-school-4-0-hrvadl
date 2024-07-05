package rate

import (
	"context"
	"fmt"

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
	return nil, nil
}

func (r *Repository) Replace(ctx context.Context, rate Exchange) error {
	if _, err := r.db.Collection(collection).InsertOne(ctx, rate); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}
