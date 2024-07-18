package cfg

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

const operation = "config parsing"

// Config struct represents application config,
// which is used application-wide.
type Config struct {
	MailerPassword      string `env:"MAILER_SMTP_PASSWORD,required,notEmpty"`
	MailerFallbackToken string `env:"MAILER_FALLBACK_TOKEN,required,notEmpty"`

	MailerFrom         string `env:"MAILER_SMTP_FROM,required,notEmpty"`
	MailerFromFallback string `env:"MAILER_FALLBACK_FROM,required,notEmpty"`

	MailerHost string `env:"MAILER_SMTP_HOST,required,notEmpty"`
	MailerPort int    `env:"MAILER_SMTP_PORT,required,notEmpty"`

	LogLevel string `env:"MAILER_LOG_LEVEL,required,notEmpty"`
	Port     string `env:"MAILER_PORT,required,notEmpty"`
	Host     string `env:"MAILER_HOST"`

	NatsURL  string `env:"NATS_URL,required,notEmpty"`
	MongoURL string `env:"MONGO_URL,required,notEmpty"`

	ConnectTimeout time.Duration `env:"MAILER_TIMEOUT"`
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
