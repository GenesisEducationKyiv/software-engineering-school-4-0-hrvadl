package mailer

import (
	"context"
	"log/slog"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v2/mailer"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

const (
	subject = "mail"
	event   = "sendDailySubMail"
)

func NewClient(pub Publisher, log *slog.Logger) *Client {
	return &Client{
		pub: pub,
		log: log,
	}
}

//go:generate mockgen -destination=./mocks/mock_publisher.go -package=mocks . Publisher
type Publisher interface {
	Publish(name string, data []byte) error
}

type Client struct {
	log *slog.Logger
	pub Publisher
}

func (c *Client) Send(ctx context.Context, html, subject string, to ...string) error {
	done := c.send(&pb.MailCommand{
		EventID:   uuid.New().String(),
		EventType: event,
		Data: &pb.Mail{
			Html:    html,
			Subject: subject,
			To:      to,
		},
	})

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (c *Client) send(m *pb.MailCommand) <-chan error {
	ch := make(chan error)

	go func() {
		bytes, err := proto.Marshal(m)
		if err != nil {
			ch <- err
			return
		}

		ch <- c.pub.Publish(subject, bytes)
	}()

	return ch
}
