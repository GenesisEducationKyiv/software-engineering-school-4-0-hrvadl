package cfg

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

const operation = "config parsing"

type Config struct {
	GatewayAddr      string `env:"GW_ADDR,required,notEmpty"`
	PromExporterAddr string `env:"PROM_EXPORTER_ADDR,required,notEmpty"`
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
	if err := env.ParseWithOptions(&cfg, env.Options{Prefix: "VEGETA_"}); err != nil {
		return nil, fmt.Errorf("%s failed: %w", operation, err)
	}
	return &cfg, nil
}
