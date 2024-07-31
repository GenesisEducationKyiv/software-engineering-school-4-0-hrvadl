package ratewatcher

import (
	"context"
	"time"
)

func NewCronJobAdapter(c Converter, t time.Duration) *Adapter {
	return &Adapter{
		timeout:  t,
		conveter: c,
	}
}

//go:generate mockgen -destination=./mocks/mock_converter.go -package=mocks . Converter
type Converter interface {
	Convert(ctx context.Context) error
}

type Adapter struct {
	timeout  time.Duration
	conveter Converter
}

func (a *Adapter) Do() error {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()
	return a.conveter.Convert(ctx)
}
