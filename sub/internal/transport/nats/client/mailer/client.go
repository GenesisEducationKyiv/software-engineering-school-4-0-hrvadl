package mailer

import (
	"context"
	"encoding/json"
	"log/slog"
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

type Mail struct {
	HTML    string   `json:"html"`
	Subject string   `json:"subject"`
	To      []string `json:"to"`
}

func (c *Client) Send(ctx context.Context, html, subject string, to ...string) error {
	done := c.send(Mail{
		HTML:    html,
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

func (c *Client) send(m Mail) <-chan error {
	ch := make(chan error)

	go func() {
		bytes, err := json.Marshal(m)
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
