package mailer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/models/mail"
)

const (
	operation = "mailer server"
	subject   = "mail"
)

func New(conn *nats.Conn, s Sender, log *slog.Logger, timeout time.Duration) *Server {
	return &Server{
		conn:    conn,
		sender:  s,
		log:     log,
		timeout: timeout,
	}
}

type Sender interface {
	Send(ctx context.Context, m mail.Mail) error
}

type Server struct {
	conn    *nats.Conn
	sender  Sender
	timeout time.Duration
	log     *slog.Logger
}

func (s *Server) Subscribe() error {
	_, err := s.conn.Subscribe(subject, s.subscribe)
	if err != nil {
		return fmt.Errorf("%s: failed to subscribe to nats: %w", operation, err)
	}
	return nil
}

func (s *Server) subscribe(msg *nats.Msg) {
	s.log.Info("Got message from NATS", slog.Any("msg", msg))
	var mail mail.Mail
	if err := json.Unmarshal(msg.Data, &mail); err != nil {
		s.log.Error("Failed to parse mail", slog.Any("err", err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if err := s.sender.Send(ctx, mail); err != nil {
		s.log.Error("Failed to send mail", slog.Any("err", err))
	}
}
