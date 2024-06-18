//go:build !integration

package cfg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	exchangeServiceBaseURLEnvKey         = "EXCHANGE_API_BASE_URL"
	exchangeFallbackServiceBaseURLEnvKey = "EXCHANGE_API_FALLBACK_BASE_URL"
	logLevelEnvKey                       = "EXCHANGE_LOG_LEVEL"
	portEnvKey                           = "EXCHANGE_PORT"
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
			},
			want: &Config{
				LogLevel:                       "debug",
				Port:                           "80",
				ExchangeServiceBaseURL:         "http://exchange.com",
				ExchangeFallbackServiceBaseURL: "http://exchange1.com",
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
				t.Setenv(exchangeServiceBaseURLEnvKey, "")
				t.Setenv(exchangeFallbackServiceBaseURLEnvKey, "http://exchange1.com")
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
