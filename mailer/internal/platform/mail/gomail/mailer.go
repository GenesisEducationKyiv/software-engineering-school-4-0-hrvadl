package gomail

import (
	"context"
	"fmt"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/mailer"
	"gopkg.in/gomail.v2"
)

const operation = "smtp client"

func NewClient(from, password, host string, port int) *Client {
	d := gomail.NewDialer(host, port, from, password)
	return &Client{
		dialer: d,
		from:   from,
	}
}

//go:generate mockgen -destination=./mocks/mock_sender.go -package=mocks . ChainedSender
type ChainedSender interface {
	Send(ctx context.Context, in *pb.Mail) error
}

type Client struct {
	dialer *gomail.Dialer
	from   string
	next   ChainedSender
}

func (c *Client) Send(ctx context.Context, in *pb.Mail) error {
	done := c.send(in)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return c.handleDone(ctx, err, in)
	}
}

func (c *Client) SetNext(next ChainedSender) {
	c.next = next
}

func (c *Client) handleDone(ctx context.Context, err error, in *pb.Mail) error {
	if err == nil {
		return nil
	}

	if c.next == nil {
		return err
	}

	if chainedErr := c.next.Send(ctx, in); chainedErr != nil {
		return fmt.Errorf("%w: %w", err, chainedErr)
	}

	return nil
}

func (c *Client) send(in *pb.Mail) <-chan error {
	done := make(chan error, 1)

	m := gomail.NewMessage()
	m.SetHeader("From", c.from)
	m.SetHeader("To", in.GetTo()...)
	m.SetHeader("Subject", in.GetSubject())
	m.SetBody("text/html", in.GetHtml())

	go func() {
		if err := c.dialer.DialAndSend(m); err != nil {
			done <- fmt.Errorf("%s: failed to dial and send: %w", operation, err)
			return
		}
		done <- nil
	}()

	return done
}
