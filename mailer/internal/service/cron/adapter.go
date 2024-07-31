package cron

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/mail"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/rate"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/subscriber"
)

const (
	operation   = "cron adapter"
	sendTimeout = time.Second * 10
)

func NewAdapter(
	rg RateGetter,
	sg SubscribersGetter,
	s Sender,
	f Formatter,
	timeout time.Duration,
	log *slog.Logger,
) *Adapter {
	return &Adapter{
		rate:        rg,
		subscribers: sg,
		sender:      s,
		timeout:     timeout,
		formatter:   f,
		log:         log,
	}
}

//go:generate mockgen -destination=./mocks/mock_fmt.go -package=mocks . Formatter
type Formatter interface {
	Format(rate float32) string
}

//go:generate mockgen -destination=./mocks/mock_rategetter.go -package=mocks . RateGetter
type RateGetter interface {
	Get(ctx context.Context) (*rate.Exchange, error)
}

//go:generate mockgen -destination=./mocks/mock_subgetter.go -package=mocks . SubscribersGetter
type SubscribersGetter interface {
	GetAll(ctx context.Context) ([]subscriber.Subscriber, error)
}

//go:generate mockgen -destination=./mocks/mock_sender.go -package=mocks . Sender
type Sender interface {
	Send(ctx context.Context, mail mail.Mail) error
}

type Adapter struct {
	rate        RateGetter
	subscribers SubscribersGetter
	sender      Sender
	formatter   Formatter
	timeout     time.Duration
	log         *slog.Logger
}

func (a *Adapter) Do() error {
	a.log.Info("Sending mails...")
	ctx, cancel := context.WithTimeout(context.Background(), sendTimeout)
	defer cancel()

	rate, err := a.rate.Get(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to get rate: %w", operation, err)
	}

	sub, err := a.subscribers.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to get subs: %w", operation, err)
	}

	m := mail.Mail{
		To:      extractMailFromSub(sub),
		Subject: "Rate exchange",
		HTML:    a.formatter.Format(rate.Rate),
	}

	if err := a.sender.Send(ctx, m); err != nil {
		return fmt.Errorf("%s: failed to send mails: %w", operation, err)
	}

	return nil
}

func extractMailFromSub(sub []subscriber.Subscriber) []string {
	mails := make([]string, 0, len(sub))
	for _, ss := range sub {
		mails = append(mails, ss.Email)
	}
	return mails
}
