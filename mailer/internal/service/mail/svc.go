package mail

import (
	"context"

	model "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/mail"
)

// NewService constructs service with provided
// default mailer.
func NewService(m Mailer) *Service {
	return &Service{
		mailers: []Mailer{m},
	}
}

//go:generate mockgen -destination=./mocks/mock_mailer.go -package=mocks . Mailer
type Mailer interface {
	Send(ctx context.Context, m model.Mail) error
}

// Service struct is responsible for aggregating and
// invoking underlying specific mailer implementations.
// If first implementation fails
// then it will call next one, until it reaches end of
// mailers slice.
type Service struct {
	mailers []Mailer
}

// Send method is responsible for invoking underlying
// specific mailer implementations. If first implementation fails
// then it will call next one, until it reaches end of
// mailers slice.
func (s *Service) Send(ctx context.Context, mail model.Mail) error {
	if mail.HTML == "" {
		return ErrEmptyContent
	}

	if mail.Subject == "" {
		return ErrEmptySubject
	}

	if len(mail.To) == 0 {
		return ErrEmptyReceivers
	}

	var err error
	for _, m := range s.mailers {
		if err = m.Send(ctx, mail); err == nil {
			return nil
		}
	}
	return err
}

// SetNext function sets new mailers to the chain of
// responsibility. Mailers are appended at the end of
// the queue.
func (s *Service) SetNext(mailers ...Mailer) {
	s.mailers = append(s.mailers, mailers...)
}
