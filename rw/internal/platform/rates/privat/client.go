package privat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const operation = "usd to uah rate"

const (
	usd = "USD"
	uah = "UAH"
)

// NewClient initializes new Client with parameters provided.
// NOTE: neither of arguments can't be empty, because in that case
// client will inevitably fail in the future.
func NewClient(url string) *Client {
	return &Client{
		url: url,
	}
}

// Client struct represents privat bank API client.
// Note: url should be a base url for  the service, not full url.
type Client struct {
	url string
}

// rate represents exchange rate API's response
// Sale is how much 1 USD is worth in a UAH.
type rate struct {
	CCY     string  `json:"ccy,omitempty"`
	BaseCCY string  `json:"base_ccy,omitempty"`
	Buy     float32 `json:"buy,omitempty,string"`
	Sale    float32 `json:"sale,omitempty,string"`
}

// Convert method converts 1 USD to UAH accordingly to the
// latest exchange rate. It's a handly wrapper around internal
// getRate() function.
func (c *Client) Convert(ctx context.Context) (float32, error) {
	res, err := c.getRate(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}
	return res, nil
}

func (c *Client) getRate(ctx context.Context) (float32, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.url+"/p24api/pubinfo?json&exchange",
		nil,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read body bytes: %w", err)
	}

	var response []rate
	if err = json.Unmarshal(bytes, &response); err != nil {
		return 0, fmt.Errorf("failed to parse response body: %w", err)
	}

	rate, err := findExchangeRateFor(uah, usd, response)
	if err != nil {
		return 0, err
	}

	return rate.Buy, nil
}

func findExchangeRateFor(base, target string, r []rate) (*rate, error) {
	for _, rr := range r {
		if rr.BaseCCY == base && rr.CCY == target {
			return &rr, nil
		}
	}
	return nil, fmt.Errorf("exchange rate for %s->%s pair is not found", base, target)
}
