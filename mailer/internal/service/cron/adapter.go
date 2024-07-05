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

const operation = "cron adapter"

func NewAdapter(
	rg RateGetter,
	sg SubscribersGetter,
	s Sender,
	timeout time.Duration,
) *Adapter {
	return &Adapter{
		rate:        rg,
		subscribers: sg,
		sender:      s,
		timeout:     timeout,
	}
}

type RateGetter interface {
	Get(context.Context) (*rate.Exchange, error)
}

type SubscribersGetter interface {
	GetAll(ctx context.Context) ([]subscriber.Subscriber, error)
}

type Sender interface {
	Send(ctx context.Context, mail mail.Mail) error
}

type Adapter struct {
	rate        RateGetter
	subscribers SubscribersGetter
	sender      Sender
	timeout     time.Duration
	log         *slog.Logger
}

func (a *Adapter) Do() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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
		HTML:    fmt.Sprintf("%v", rate),
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
