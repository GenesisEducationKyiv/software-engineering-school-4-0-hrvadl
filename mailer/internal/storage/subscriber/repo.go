package subscriber

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/platform/db"
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
	cur, err := r.db.Collection(collection).Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get collection: %w", operation, err)
	}

	var sub []Subscriber
	err = cur.All(ctx, &sub)
	if err == nil {
		return sub, nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("%s: %w: %w", operation, db.ErrNotFound, err)
	}

	return nil, fmt.Errorf("%s: failed to decode subs: %w", operation, err)
}

func (r *Repository) Save(ctx context.Context, sub Subscriber) error {
	if _, err := r.db.Collection(collection).InsertOne(ctx, sub); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, sub Subscriber) error {
	res, err := r.db.Collection(collection).DeleteOne(ctx, sub)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	if res.DeletedCount == 0 {
		return db.ErrNotFound
	}

	return nil
}
