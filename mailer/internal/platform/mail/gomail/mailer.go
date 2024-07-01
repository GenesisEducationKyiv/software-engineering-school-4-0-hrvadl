package gomail

import (
	"context"
	"fmt"

	"gopkg.in/gomail.v2"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/models/mail"
)

const operation = "smtp client"

// NewClient constructs new Gomail client
// with provided credentials.
func NewClient(from, password, host string, port int) *Client {
	d := gomail.NewDialer(host, port, from, password)
	return &Client{
		dialer: d,
		from:   from,
	}
}

// Client is a thin wrapper around Gomail library
// which will add context support to the existing
// signature call.
type Client struct {
	dialer *gomail.Dialer
	from   string
}

// Send method initiates a call to the SMTP server using
// Gomail's method. Blocks until call is finished, or
// error is raised, or context is done.
func (c *Client) Send(ctx context.Context, in mail.Mail) error {
	done := c.send(in)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (c *Client) send(in mail.Mail) <-chan error {
	done := make(chan error, 1)

	go func() {
		m := gomail.NewMessage()
		m.SetHeader("From", c.from)
		m.SetHeader("Bcc", in.To...)
		m.SetHeader("Subject", in.Subject)
		m.SetBody("text/html", in.HTML)

		if err := c.dialer.DialAndSend(m); err != nil {
			done <- fmt.Errorf("%s: failed to dial and send: %w", operation, err)
			return
		}

		done <- nil
	}()

	return done
}
