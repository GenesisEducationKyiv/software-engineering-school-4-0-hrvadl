package nats

import "github.com/nats-io/nats.go"

func Must(c *Client, err error) *Client {
	if err != nil {
		panic(err)
	}
	return c
}

func NewClient(url string) (*Client, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

type Client struct {
	conn *nats.Conn
}

func (c *Client) Publish(name string, data []byte) error {
	return c.conn.Publish(name, data)
}
