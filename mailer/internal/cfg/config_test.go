//go:build !integration

package cfg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	logLevelEnvKey            = "MAILER_LOG_LEVEL"
	portEnvKey                = "MAILER_PORT"
	mailerTokenEnvKey         = "MAILER_SMTP_PASSWORD"  // #nosec G101
	mailerFallbackTokenEnvKey = "MAILER_FALLBACK_TOKEN" // #nosec G101
	mailerFromEnvKey          = "MAILER_SMTP_FROM"
	mailerFallbackFromEnvKey  = "MAILER_FALLBACK_FROM"
	mailerHostEnvKey          = "MAILER_SMTP_HOST"
	mailerPortEnvKey          = "MAILER_SMTP_PORT"
	natsURLEnvKey             = "NATS_URL"
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
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "80")
				t.Setenv(mailerTokenEnvKey, "secret")
				t.Setenv(mailerFallbackTokenEnvKey, "secret")
				t.Setenv(mailerFallbackFromEnvKey, "secret@test.com")
				t.Setenv(mailerFromEnvKey, "secret@test.com")
				t.Setenv(mailerHostEnvKey, "smtp.google.com")
				t.Setenv(mailerPortEnvKey, "528")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want: &Config{
				LogLevel:            "debug",
				Port:                "80",
				MailerPassword:      "secret",
				MailerFallbackToken: "secret",
				MailerFrom:          "secret@test.com",
				MailerFromFallback:  "secret@test.com",
				MailerHost:          "smtp.google.com",
				MailerPort:          528,
				NatsURL:             "nats://nats:4222",
			},
			wantErr: false,
		},
		{
			name: "Should not parse config correctly when log level is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(portEnvKey, "80")
				t.Setenv(mailerTokenEnvKey, "secret")
				t.Setenv(mailerFallbackTokenEnvKey, "secret")
				t.Setenv(mailerFromEnvKey, "secret@test.com")
				t.Setenv(mailerFallbackFromEnvKey, "secret@test.com")
				t.Setenv(mailerHostEnvKey, "smtp.google.com")
				t.Setenv(mailerPortEnvKey, "528")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config correctly when smpt port is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(mailerTokenEnvKey, "secret")
				t.Setenv(mailerFallbackTokenEnvKey, "secret")
				t.Setenv(mailerFromEnvKey, "secret@test.com")
				t.Setenv(mailerFallbackFromEnvKey, "secret@test.com")
				t.Setenv(mailerHostEnvKey, "smtp.google.com")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config correctly when password is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "80")
				t.Setenv(mailerFallbackTokenEnvKey, "secret")
				t.Setenv(mailerFromEnvKey, "secret@test.com")
				t.Setenv(mailerFallbackFromEnvKey, "secret@test.com")
				t.Setenv(mailerHostEnvKey, "smtp.google.com")
				t.Setenv(mailerPortEnvKey, "528")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config correctly when token is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "80")
				t.Setenv(mailerTokenEnvKey, "secret")
				t.Setenv(mailerFromEnvKey, "secret@test.com")
				t.Setenv(mailerFallbackFromEnvKey, "secret@test.com")
				t.Setenv(mailerHostEnvKey, "smtp.google.com")
				t.Setenv(mailerPortEnvKey, "528")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when from mail is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "80")
				t.Setenv(mailerTokenEnvKey, "secret")
				t.Setenv(mailerFallbackTokenEnvKey, "secret")
				t.Setenv(mailerFallbackFromEnvKey, "secret@test.com")
				t.Setenv(mailerHostEnvKey, "smtp.google.com")
				t.Setenv(mailerPortEnvKey, "528")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when fallback from mail is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "80")
				t.Setenv(mailerTokenEnvKey, "secret")
				t.Setenv(mailerFallbackTokenEnvKey, "secret")
				t.Setenv(mailerFromEnvKey, "secret@test.com")
				t.Setenv(mailerHostEnvKey, "smtp.google.com")
				t.Setenv(mailerPortEnvKey, "528")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when smtp port is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "80")
				t.Setenv(mailerTokenEnvKey, "secret")
				t.Setenv(mailerFallbackTokenEnvKey, "secret")
				t.Setenv(mailerFallbackFromEnvKey, "secret@test.com")
				t.Setenv(mailerFromEnvKey, "secret@test.com")
				t.Setenv(mailerHostEnvKey, "host")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when smtp host is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "80")
				t.Setenv(mailerTokenEnvKey, "secret")
				t.Setenv(mailerFallbackTokenEnvKey, "secret")
				t.Setenv(mailerFallbackFromEnvKey, "secret@test.com")
				t.Setenv(mailerFromEnvKey, "secret@test.com")
				t.Setenv(mailerPortEnvKey, "528")
				t.Setenv(natsURLEnvKey, "nats://nats:4222")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should not parse config when nats url is missing",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv(logLevelEnvKey, "debug")
				t.Setenv(portEnvKey, "80")
				t.Setenv(mailerTokenEnvKey, "secret")
				t.Setenv(mailerFallbackTokenEnvKey, "secret")
				t.Setenv(mailerFallbackFromEnvKey, "secret@test.com")
				t.Setenv(mailerFromEnvKey, "secret@test.com")
				t.Setenv(mailerPortEnvKey, "528")
				t.Setenv(mailerHostEnvKey, "smtp.google.com")
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
