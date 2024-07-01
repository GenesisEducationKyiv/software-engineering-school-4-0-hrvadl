package mailer

import (
	"context"
	"log/slog"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/mailer"
	"google.golang.org/protobuf/proto"
)

const subject = "mail"

func NewClient(pub Publisher, log *slog.Logger) *Client {
	return &Client{
		pub: pub,
		log: log,
	}
}

type Publisher interface {
	Publish(name string, data []byte) error
}

type Client struct {
	log *slog.Logger
	pub Publisher
}

func (c *Client) Send(ctx context.Context, html, subject string, to ...string) error {
	done := c.send(&pb.Mail{
		Html:    html,
		Subject: subject,
		To:      to,
	})

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (c *Client) send(m *pb.Mail) <-chan error {
	ch := make(chan error)

	go func() {
		bytes, err := proto.Marshal(m)
		if err != nil {
			ch <- err
			return
		}

		if err := c.pub.Publish(subject, bytes); err != nil {
			ch <- err
			return
		}
		ch <- nil
	}()

	return ch
}
