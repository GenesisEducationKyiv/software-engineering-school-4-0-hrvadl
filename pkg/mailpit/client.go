package mailpit

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Receipient struct {
	Address string `json:"Address,omitempty"`
	Name    string `json:"Name,omitempty"`
}

type Message struct {
	ID      string       `json:"ID,omitempty"`
	From    Receipient   `json:"From,omitempty"`
	To      []Receipient `json:"To,omitempty"`
	Bcc     []Receipient `json:"Bcc,omitempty"`
	Cc      []Receipient `json:"Cc,omitempty"`
	Subject string       `json:"Subject,omitempty"`
}

type Response struct {
	Messages []Message `json:"messages,omitempty"`
}

func NewClient(host string, port int, timeout time.Duration) *Client {
	return &Client{
		addr: net.JoinHostPort(host, strconv.Itoa(port)),
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

type Client struct {
	addr   string
	client *http.Client
}

func (c *Client) GetAll() ([]Message, error) {
	r, err := http.NewRequest(http.MethodGet, c.toURL("/api/v1/messages"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create req: %w", err)
	}

	res, err := c.client.Do(r)
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

	var msg Response
	if err := json.Unmarshal(b, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshall: %w", err)
	}

	return msg.Messages, nil
}

func (c *Client) DeleteAll() error {
	r, err := http.NewRequest(http.MethodDelete, c.toURL("/api/v1/messages"), nil)
	if err != nil {
		return fmt.Errorf("failed to create req: %w", err)
	}

	res, err := c.client.Do(r)
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
