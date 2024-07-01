package cfg

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

const operation = "config parsing"

// Config struct represents application config,
// which is used application-wide.
type Config struct {
	MailerAddr      string `env:"MAILER_ADDR,required,notEmpty"`
	Dsn             string `env:"SUB_DSN,required,notEmpty"`
	RateWatcherAddr string `env:"RATE_WATCH_ADDR,required,notEmpty"`
	Port            string `env:"SUB_PORT,required,notEmpty"`
	LogLevel        string `env:"SUB_LOG_LEVEL,required,notEmpty"`
	NatsURL         string `env:"NATS_URL,required,notEmpty"`
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
