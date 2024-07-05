package subscriber

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collection = "subscribers"
	operation  = "subscribers repo"
)

func NewRepo(db *mongo.Database) *Repository {
	return &Repository{
		db: db,
	}
}

type Repository struct {
	db *mongo.Database
}

func (r *Repository) GetAll(ctx context.Context) ([]Subscriber, error) {
	return nil, nil
}

func (r *Repository) Save(ctx context.Context, sub Subscriber) error {
	slog.Info("Saving subscriber", slog.Any("sub", sub))
	if _, err := r.db.Collection(collection).InsertOne(ctx, sub); err != nil {
		slog.Error("Failed to save sub", slog.Any("err", err))
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, sub Subscriber) error {
	if _, err := r.db.Collection(collection).DeleteOne(ctx, sub); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}
