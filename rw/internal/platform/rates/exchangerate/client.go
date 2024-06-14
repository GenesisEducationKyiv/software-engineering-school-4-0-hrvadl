package exchangerate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	operation = "usd to uah rate"
)

const (
	usd = "USD"
	uah = "UAH"
)

// usdUahResponse represents exchange rate API's response
// ConversionRate is how much 1 USD is worth in a UAH.
// Struct also contains some meta fields which can be useful in long run,
// such as UpdatedAt and TargetCode.
type usdUahResponse struct {
	ConversionRate float32 `json:"conversion_rate"`
	BaseCode       string  `json:"base_code"`
	TargetCode     string  `json:"target_code"`
	UpdatedAt      int     `json:"time_last_update_unix"`
}

// NewClient initializes new Client with parameters provided.
// NOTE: neither of arguments can't be empty, because in that case
// client will inevitably fail in the future.
func NewClient(token, url string) Client {
	return Client{
		token: token,
		url:   url,
	}
}

// Client struct represents exchange rate API client.
// Note: url should be a base url for  the service, not full url.
type Client struct {
	token string
	url   string
}

// Convert method converts 1 USD to UAH accordingly to the
// latest exchange rate. It's a handly wrapper around internal
// getRate() function.
func (c Client) Convert(ctx context.Context) (float32, error) {
	res := new(usdUahResponse)
	if err := c.getRate(ctx, res, usd, uah); err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	return res.ConversionRate, nil
}

// getRate method is used to query how much **from** currency is worth
// in **to**  currency. response should be a pointer to the API response.
// It's needed because API returns different responses for differect currency pairs.
func (c Client) getRate(
	ctx context.Context,
	response any,
	from string,
	to string,
) error {
	if c.token == "" {
		return errors.New("token is missing")
	}

	url, err := url.Parse(fmt.Sprintf("%s/%s/pair/%s/%s", c.url, c.token, from, to))
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to construct request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read body bytes: %w", err)
	}

	if err := json.Unmarshal(bytes, &response); err != nil {
		return fmt.Errorf("failed to parse response body: %w", err)
	}

	return nil
}
