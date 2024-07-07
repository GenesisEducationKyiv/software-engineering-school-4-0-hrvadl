package rate

import (
	"context"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/platform/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/rate"
)

func NewService(rs RateSource) *Service {
	return &Service{
		rate: rs,
	}
}

type RateGetter interface {
	Get(ctx context.Context) (*rate.Exchange, error)
}

type RateReplacer interface {
	Replace(ctx context.Context, rate rate.Exchange) error
}

//go:generate mockgen -destination=./mocks/mock_ratesource.go -package=mocks . RateSource
type RateSource interface {
	RateGetter
	RateReplacer
}

type Service struct {
	rate RateSource
}

func (s *Service) Get(ctx context.Context) (*rate.Exchange, error) {
	rate, err := s.rate.Get(ctx)
	if err == nil {
		return rate, nil
	}

	if errors.Is(err, db.ErrNotFound) {
		return nil, ErrNotFound
	}

	return nil, fmt.Errorf("%w: %w", ErrFailetToGet, err)
}

func (s *Service) Replace(ctx context.Context, rate rate.Exchange) error {
	if rate.Rate == 0 {
		return ErrEmptyRate
	}

	if err := s.rate.Replace(ctx, rate); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToReplace, err)
	}

	return nil
}
