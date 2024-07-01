package rw

import "context"

// NewService constructs service with provided
// default exchange rate source.
func NewService(source RateSource) *Service {
	return &Service{
		sources: []RateSource{source},
	}
}

//go:generate mockgen -destination=./mocks/mock_source.go -package=mocks . RateSource
type RateSource interface {
	Convert(ctx context.Context) (float32, error)
}

// Service struct is responsible for aggregating and
// invoking underlying specific rate source implementations.
// If first implementation fails
// then it will call next one, until it reaches end of
// mailers slice.
type Service struct {
	sources []RateSource
}

// Convert method is responsible for invoking underlying
// specific rate source implementations. If first implementation fails
// then it will call next one, until it reaches end of
// mailers slice.
func (s *Service) Convert(ctx context.Context) (float32, error) {
	var (
		rate float32
		err  error
	)

	for _, svc := range s.sources {
		rate, err = svc.Convert(ctx)
		if err == nil {
			return rate, nil
		}
	}

	return rate, err
}

// SetNext function sets new sources to the chain of
// responsibility. Rate sources are appended at
// the end of the queue.
func (s *Service) SetNext(source ...RateSource) {
	s.sources = append(s.sources, source...)
}
