package mail

import (
	"context"

	model "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/models/mail"
)

func NewService(m Mailer) *Service {
	return &Service{
		mailers: []Mailer{m},
	}
}

type Mailer interface {
	Send(ctx context.Context, m model.Mail) error
}

type Service struct {
	mailers []Mailer
}

func (s *Service) Send(ctx context.Context, mail model.Mail) error {
	var err error
	for _, m := range s.mailers {
		if err = m.Send(ctx, mail); err == nil {
			return nil
		}
	}
	return err
}

func (s *Service) SetNext(mailers ...Mailer) {
	s.mailers = append(s.mailers, mailers...)
}
