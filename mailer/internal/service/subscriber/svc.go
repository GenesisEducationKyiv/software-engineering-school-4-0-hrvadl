package subscriber

import (
	"context"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/platform/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/subscriber"
)

func NewService(ss SubscriberSource) *Service {
	return &Service{
		subscriber: ss,
	}
}

type SubscriberSaver interface {
	Save(ctx context.Context, sub subscriber.Subscriber) error
}

type SubscriberDeleter interface {
	Delete(ctx context.Context, sub subscriber.Subscriber) error
}

type SubscriberGetter interface {
	GetAll(ctx context.Context) ([]subscriber.Subscriber, error)
}

//go:generate mockgen -destination=./mocks/mock_subsource.go -package=mocks . SubscriberSource
type SubscriberSource interface {
	SubscriberSaver
	SubscriberGetter
	SubscriberDeleter
}

type Service struct {
	subscriber SubscriberSource
}

func (s *Service) GetAll(ctx context.Context) ([]subscriber.Subscriber, error) {
	sub, err := s.subscriber.GetAll(ctx)
	if err == nil {
		return sub, nil
	}

	if errors.Is(err, db.ErrNotFound) {
		return nil, ErrNoSubscribers
	}

	return nil, fmt.Errorf("%w: %w", ErrFailedToGet, err)
}

func (s *Service) Save(ctx context.Context, sub subscriber.Subscriber) error {
	if err := s.subscriber.Save(ctx, sub); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSave, err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, sub subscriber.Subscriber) error {
	if err := s.subscriber.Delete(ctx, sub); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToDelete, err)
	}
	return nil
}
