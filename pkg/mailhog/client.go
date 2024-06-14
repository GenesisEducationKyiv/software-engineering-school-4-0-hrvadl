package mailhog

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"
)

type From struct {
	Domain  string `json:"domain,omitempty"`
	Mailbox string `json:"mailbox,omitempty"`
}

type To struct {
	Domain  string `json:"domain,omitempty"`
	Mailbox string `json:"mailbox,omitempty"`
}

type Headers struct {
	From    []string `json:"from,omitempty"`
	To      []string `json:"to,omitempty"`
	Subject []string `json:"subject,omitempty"`
}

type Content struct {
	Body    string  `json:"body,omitempty"`
	Headers Headers `json:"headers,omitempty"`
}

type Mail struct {
	ID      string  `json:"id,omitempty"`
	From    From    `json:"from,omitempty"`
	To      []To    `json:"to,omitempty"`
	Content Content `json:"content,omitempty"`
}

func NewClient(host string, port int, timeout time.Duration) *Client {
	return &Client{
		addr: net.JoinHostPort(host, strconv.Itoa(port)),
		cl: &http.Client{
			Timeout: timeout,
		},
	}
}

type Client struct {
	addr string
	cl   *http.Client
}

func (c *Client) GetAll() ([]Mail, error) {
	r, err := http.NewRequest(http.MethodGet, c.toURL("/api/v1/messages"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create req: %w", err)
	}

	res, err := c.cl.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to send req: %w", err)
	}

	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mailhog returned negative status code: %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes: %w", err)
	}

	var msg []Mail
	slog.Info("Received response", slog.Any("bytes", b))
	if err := json.Unmarshal(b, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshall: %w", err)
	}

	return msg, nil
}

func (c *Client) DeleteAll() error {
	r, err := http.NewRequest(http.MethodDelete, c.toURL("/api/v1/messages"), nil)
	if err != nil {
		return fmt.Errorf("failed to create req: %w", err)
	}

	res, err := c.cl.Do(r)
	if err != nil {
		return fmt.Errorf("failed to send req: %w", err)
	}

	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("mailhog returned negative status code: %d", res.StatusCode)
	}

	return nil
}

func (c *Client) toURL(h string) string {
	return "http://" + c.addr + h
}
