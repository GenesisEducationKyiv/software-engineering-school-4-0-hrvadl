//go:build !integration

package cfg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	mailerServiceAddrEnvKey = "MAILER_ADDR"
	rateWatchAddrEnvKey     = "RATE_WATCH_ADDR"
	logLevelEnvKey          = "SUB_LOG_LEVEL"
	portEnvKey              = "SUB_PORT"
	dsnEnvKey               = "SUB_DSN"
	natsURLEnvKey           = "NATS_URL"
	prometheusPortEnvKey    = "PROMETHEUS_PORT"
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
					Must(tt.args.cfg, tt.args.err)
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
				t.Setenv(mailerServiceAddrEnvKey, "mailer:80")
				t.Setenv(rateWatchAddrEnvKey, "rw:8080")
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "3030")
				t.Setenv(dsnEnvKey, "mysql://test:tests@(db:testse)/shgsoh")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
				t.Setenv(prometheusPortEnvKey, "2112")
			},
			want: &Config{
				MailerAddr:      "mailer:80",
				RateWatcherAddr: "rw:8080",
				LogLevel:        "debug",
				Port:            "3030",
				Dsn:             "mysql://test:tests@(db:testse)/shgsoh",
				NatsURL:         "nats://nats:4222",
				PrometheusPort:  "2112",
			},
			wantErr: false,
		},
		{
			name: "Should not parse config when mailer addr is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(rateWatchAddrEnvKey, "rw:8080")
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "3030")
				t.Setenv(dsnEnvKey, "mysql://test:tests@(db:testse)/shgsoh")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
				t.Setenv(prometheusPortEnvKey, "2112")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when rw addr is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(mailerServiceAddrEnvKey, "mailer:80")
				t.Setenv(rateWatchAddrEnvKey, "")
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "3030")
				t.Setenv(dsnEnvKey, "mysql://test:tests@(db:testse)/shgsoh")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
				t.Setenv(prometheusPortEnvKey, "2112")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when log level is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(mailerServiceAddrEnvKey, "mailer:80")
				t.Setenv(rateWatchAddrEnvKey, "rw:2209")
				t.Setenv(logLevelEnvKey, "")
				t.Setenv(portEnvKey, "3030")
				t.Setenv(dsnEnvKey, "mysql://test:tests@(db:testse)/shgsoh")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
				t.Setenv(prometheusPortEnvKey, "2112")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when port is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(mailerServiceAddrEnvKey, "mailer:80")
				t.Setenv(rateWatchAddrEnvKey, "rw:2209")
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "")
				t.Setenv(dsnEnvKey, "mysql://test:tests@(db:testse)/shgsoh")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
				t.Setenv(prometheusPortEnvKey, "2112")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when dsn is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(mailerServiceAddrEnvKey, "mailer:80")
				t.Setenv(rateWatchAddrEnvKey, "rw:2209")
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "801")
				t.Setenv(dsnEnvKey, "")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
				t.Setenv(prometheusPortEnvKey, "2112")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when nats URL is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(mailerServiceAddrEnvKey, "mailer:80")
				t.Setenv(rateWatchAddrEnvKey, "rw:2209")
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "4040")
				t.Setenv(dsnEnvKey, "mysql://test:tests@(db:testse)/shgsoh")
				t.Setenv(prometheusPortEnvKey, "2112")
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
