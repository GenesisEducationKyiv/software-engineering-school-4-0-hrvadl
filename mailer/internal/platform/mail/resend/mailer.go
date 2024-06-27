package resend

import (
	"context"
	"fmt"

	rs "github.com/resend/resend-go/v2"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/models/mail"
)

const operation = "resend mail client"

// NewClient constructs new Resend client
// with provided token.
func NewClient(from, token string) *Client {
	return &Client{
		client: rs.NewClient(token),
		from:   from,
	}
}

// Client is a thin wrapper around resend's SDK
// which will add context support to the existing
// signature call.
type Client struct {
	client *rs.Client
	from   string
}

// Send method initiates a call to the resend API using
// bult-in resend's SDK. Blocks until call is finished, or
// error is raised, or context is done.
func (c *Client) Send(ctx context.Context, m mail.Mail) error {
	if len(m.To) == 0 {
		return fmt.Errorf("%s: recipients cannot be empty", operation)
	}

	done := c.send(m)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (c *Client) send(m mail.Mail) <-chan error {
	done := make(chan error, 1)

	go func() {
		_, err := c.client.Emails.Send(&rs.SendEmailRequest{
			From:    c.from,
			To:      m.To,
			Subject: m.Subject,
			Html:    m.HTML,
		})
		if err != nil {
			done <- fmt.Errorf("%s: failed to send message: %w", operation, err)
			return
		}
		done <- nil
	}()

	return done
}
