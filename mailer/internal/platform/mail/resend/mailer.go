package resend

import (
	"context"
	"fmt"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/mailer"
	rs "github.com/resend/resend-go/v2"
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

type ChainedSender interface {
	Send(ctx context.Context, in *pb.Mail) error
}

// Client is a thin wrapper around resend's SDK
// which will add context support to the existing
// signature call.
type Client struct {
	client *rs.Client
	from   string
	next   ChainedSender
}

// Send method initiates a call to the resend API using
// bult-in resend's SDK. Blocks until call is finished, or
// error is raised, or context is done.
func (c *Client) Send(ctx context.Context, m *pb.Mail) error {
	if len(m.GetTo()) == 0 {
		return fmt.Errorf("%s: recipients cannot be empty", operation)
	}

	done := c.send(m)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return c.handleDone(ctx, err, m)
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

func (c *Client) send(m *pb.Mail) <-chan error {
	done := make(chan error, 1)

	go func() {
		_, err := c.client.Emails.Send(&rs.SendEmailRequest{
			From:    c.from,
			To:      m.GetTo(),
			Subject: m.GetSubject(),
			Html:    m.GetHtml(),
		})
		if err != nil {
			done <- fmt.Errorf("%s: failed to send message: %w", operation, err)
			return
		}
		done <- nil
	}()

	return done
}
