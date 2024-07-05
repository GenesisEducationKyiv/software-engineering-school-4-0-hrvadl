package cfg

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

const operation = "config parsing"

// Config struct represents application config,
// which is used application-wide.
type Config struct {
	ExchangeServiceBaseURL               string `env:"EXCHANGE_API_BASE_URL,required,notEmpty"`
	ExchangeFallbackServiceBaseURL       string `env:"EXCHANGE_API_FALLBACK_BASE_URL,required,notEmpty"`
	ExchangeFallbackSecondServiceBaseURL string `env:"EXCHANGE_API_FALLBACK2_BASE_URL,required,notEmpty"`
	ExchangeFallbackSecondServiceToken   string `env:"EXCHANGE_API_FALLBACK2_TOKEN,required,notEmpty"`
	LogLevel                             string `env:"EXCHANGE_LOG_LEVEL,required,notEmpty"`
	Port                                 string `env:"EXCHANGE_PORT,required,notEmpty"`
	NatsURL                              string `env:"NATS_URL,required,notEmpty"`
}

// Must is a handly wrapper around return results from
// the NewFromEnv() function, which will panic in case of error.
// Should be called only in main function, when we don't need
// to handle errors.
func Must(cfg *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return cfg
}

// NewFromEnv parses the environment variables into
// the Config struct. Returns an error if any of required variables
// is missing or contains invalid value.
func NewFromEnv() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("%s failed: %w", operation, err)
	}
	return &cfg, nil
}
