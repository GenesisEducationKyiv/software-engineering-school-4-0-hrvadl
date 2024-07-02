package mailer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v2/mailer"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

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
	var in pb.MailEvent
	if err := proto.Unmarshal(msg.Data, &in); err != nil {
		s.log.Error("Failed to parse mail", slog.Any("err", err))
		return
	}

	s.log.Info(
		"Got message from NATS",
		slog.String("id", in.GetEventID()),
		slog.String("type", in.GetEventType()),
	)

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	data := in.GetData()

	mail := mail.Mail{
		HTML:    data.GetHtml(),
		To:      data.GetTo(),
		Subject: data.GetSubject(),
	}

	if err := s.sender.Send(ctx, mail); err != nil {
		s.nack(msg)
		return
	}

	s.ack(msg)
}

func (s *Server) ack(msg *nats.Msg) {
	if err := msg.Ack(); err != nil {
		s.log.Error("Failed to send ack", slog.Any("err", err))
	}
}

func (s *Server) nack(msg *nats.Msg) {
	if err := msg.Nak(); err != nil {
		s.log.Error("Failed to send ack", slog.Any("err", err))
	}
}
