package rates

import (
	"context"
	"log/slog"
)

func NewWithLogger(base Converter, log *slog.Logger) *WithLoggerDecorator {
	return &WithLoggerDecorator{
		base: base,
		log:  log,
	}
}

//go:generate mockgen -destination=./mocks/mock_converter.go -package=mocks . Converter
type Converter interface {
	Convert(ctx context.Context) (float32, error)
}

type WithLoggerDecorator struct {
	base Converter
	log  *slog.Logger
}

func (c WithLoggerDecorator) Convert(ctx context.Context) (float32, error) {
	c.log.Info("Sending request to exchange API service")
	res, err := c.base.Convert(ctx)
	if err != nil {
		c.log.Error("Received error", slog.Any("err", err))
		return res, err
	}

	c.log.Error("Received response from API", slog.Any("res", res))
	return res, err
}
