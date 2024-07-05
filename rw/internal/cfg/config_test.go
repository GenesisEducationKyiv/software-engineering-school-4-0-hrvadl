//go:build !integration

package cfg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	exchangeServiceBaseURLEnvKey               = "EXCHANGE_API_BASE_URL"
	exchangeFallbackServiceBaseURLEnvKey       = "EXCHANGE_API_FALLBACK_BASE_URL"
	exchangeFallbackSecondServiceBaseURLEnvKey = "EXCHANGE_API_FALLBACK2_BASE_URL"
	exchangeFallbackSecondServiceToken         = "EXCHANGE_API_FALLBACK2_TOKEN" // #nosec G101
	logLevelEnvKey                             = "EXCHANGE_LOG_LEVEL"
	natsURLEnvKey                              = "NATS_URL"
	portEnvKey                                 = "EXCHANGE_PORT"
)

func TestMust(t *testing.T) {
	t.Parallel()
	type args struct {
		cfg *Config
		err error
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "Should not panic when err is nil",
			args: args{
				cfg: &Config{},
				err: nil,
			},
			want: &Config{},
		},
		{
			name: "Should panic when err is not nil",
			args: args{
				cfg: nil,
				err: errors.New("failed to parse config"),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.wantErr {
				require.Panics(t, func() {
					_ = Must(tt.args.cfg, tt.args.err)
				})
				return
			}

			require.NotPanics(t, func() {
				got := Must(tt.args.cfg, tt.args.err)
				require.Equal(t, tt.want, got)
			})
		})
	}
}

func TestNewFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T)
		want    *Config
		wantErr bool
	}{
		{
			name: "Should parse config correctly when all env vars are present",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "80")
				t.Setenv(exchangeServiceBaseURLEnvKey, "http://exchange.com")
				t.Setenv(exchangeFallbackServiceBaseURLEnvKey, "http://exchange1.com")
				t.Setenv(exchangeFallbackSecondServiceBaseURLEnvKey, "http://exchange2.com")
				t.Setenv(exchangeFallbackSecondServiceToken, "token")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want: &Config{
				LogLevel:                             "debug",
				Port:                                 "80",
				ExchangeServiceBaseURL:               "http://exchange.com",
				ExchangeFallbackServiceBaseURL:       "http://exchange1.com",
				ExchangeFallbackSecondServiceBaseURL: "http://exchange2.com",
				ExchangeFallbackSecondServiceToken:   "token",
				NatsURL:                              "nats://nats:4222",
			},
			wantErr: false,
		},
		{
			name: "Should not parse config when log level is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "")
				t.Setenv(portEnvKey, "80")
				t.Setenv(exchangeServiceBaseURLEnvKey, "http://exchange.com")
				t.Setenv(exchangeFallbackServiceBaseURLEnvKey, "http://exchange1.com")
				t.Setenv(exchangeFallbackSecondServiceBaseURLEnvKey, "http://exchange2.com")
				t.Setenv(exchangeFallbackSecondServiceToken, "token")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when port is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "")
				t.Setenv(exchangeServiceBaseURLEnvKey, "http://exchange.com")
				t.Setenv(exchangeFallbackServiceBaseURLEnvKey, "http://exchange1.com")
				t.Setenv(exchangeFallbackSecondServiceBaseURLEnvKey, "http://exchange2.com")
				t.Setenv(exchangeFallbackSecondServiceToken, "token")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when base URL is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "80")
				t.Setenv(exchangeFallbackServiceBaseURLEnvKey, "http://exchange1.com")
				t.Setenv(exchangeFallbackSecondServiceBaseURLEnvKey, "http://exchange2.com")
				t.Setenv(exchangeFallbackSecondServiceToken, "token")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when fallback base URL is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "80")
				t.Setenv(exchangeServiceBaseURLEnvKey, "http://exchange.com")
				t.Setenv(exchangeFallbackSecondServiceBaseURLEnvKey, "http://exchange2.com")
				t.Setenv(exchangeFallbackSecondServiceToken, "token")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when second fallback token is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "80")
				t.Setenv(exchangeServiceBaseURLEnvKey, "http://exchange.com")
				t.Setenv(exchangeFallbackServiceBaseURLEnvKey, "http://exchange1.com")
				t.Setenv(exchangeFallbackSecondServiceBaseURLEnvKey, "http://exchange2.com")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when second fallback base URL is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "80")
				t.Setenv(exchangeServiceBaseURLEnvKey, "http://exchange.com")
				t.Setenv(exchangeFallbackServiceBaseURLEnvKey, "http://exchange1.com")
				t.Setenv(exchangeFallbackSecondServiceToken, "token")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when second fallback base URL is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "80")
				t.Setenv(exchangeServiceBaseURLEnvKey, "http://exchange.com")
				t.Setenv(exchangeFallbackServiceBaseURLEnvKey, "http://exchange1.com")
				t.Setenv(exchangeFallbackSecondServiceToken, "token")
				t.Setenv(exchangeFallbackSecondServiceBaseURLEnvKey, "http://url.com")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			got, err := NewFromEnv()
			require.Equal(t, tt.want, got)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
