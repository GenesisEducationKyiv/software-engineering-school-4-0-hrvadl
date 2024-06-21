//go:build integration

package gomail

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/pkg/mailpit"
	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/mailer"
	"github.com/stretchr/testify/require"
)

const (
	mailerSMTPHostEnvKey     = "MAILER_TEST_SMTP_HOST"
	mailerSMTPPasswordEnvKey = "MAILER_TEST_SMTP_PASSWORD" // #nosec G101
	mailerSMTPFromEnvKey     = "MAILER_TEST_SMTP_FROM"
	mailerSMTPPortEnvKey     = "MAILER_TEST_SMTP_PORT"
	mailerTestAPIPortEnvKey  = "MAILER_TEST_API_PORT"
)

func TestClientSend(t *testing.T) {
	type fields struct {
		from     string
		password string
		host     string
		port     int
	}
	type args struct {
		ctx context.Context
		in  *pb.Mail
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Should send mail correctly",
			wantErr: false,
			args: args{
				ctx: context.Background(),
				in: &pb.Mail{
					To:      []string{"test@gmail.com"},
					Html:    "hello",
					Subject: "world",
				},
			},
			fields: fields{
				from:     mustGetEnv(t, mailerSMTPFromEnvKey),
				password: mustGetEnv(t, mailerSMTPPasswordEnvKey),
				host:     mustGetEnv(t, mailerSMTPHostEnvKey),
				port:     mustAtoi(t, mustGetEnv(t, mailerSMTPPortEnvKey)),
			},
		},
		{
			name:    "Should send mail correctly to multiple receivers",
			wantErr: false,
			args: args{
				ctx: context.Background(),
				in: &pb.Mail{
					To:      []string{"test@gmail.com", "test2@gmail.com"},
					Html:    "hello",
					Subject: "world",
				},
			},
			fields: fields{
				from:     mustGetEnv(t, mailerSMTPFromEnvKey),
				password: mustGetEnv(t, mailerSMTPPasswordEnvKey),
				host:     mustGetEnv(t, mailerSMTPHostEnvKey),
				port:     mustAtoi(t, mustGetEnv(t, mailerSMTPPortEnvKey)),
			},
		},
		{
			name:    "Should return error when it takes too long",
			wantErr: true,
			args: args{
				ctx: newImmediateCtx(),
				in: &pb.Mail{
					To:      []string{"test@gmail.com", "invalid.com"},
					Html:    "hello",
					Subject: "world",
				},
			},
			fields: fields{
				from:     mustGetEnv(t, mailerSMTPFromEnvKey),
				password: mustGetEnv(t, mailerSMTPPasswordEnvKey),
				host:     mustGetEnv(t, mailerSMTPHostEnvKey),
				port:     mustAtoi(t, mustGetEnv(t, mailerSMTPPortEnvKey)),
			},
		},
		{
			name:    "Should return error when invalid email is given",
			wantErr: true,
			args: args{
				ctx: context.Background(),
				in: &pb.Mail{
					To:      []string{"test@gmail.com", "invalid.com"},
					Html:    "hello",
					Subject: "world",
				},
			},
			fields: fields{
				from:     mustGetEnv(t, mailerSMTPFromEnvKey),
				password: mustGetEnv(t, mailerSMTPPasswordEnvKey),
				host:     mustGetEnv(t, mailerSMTPHostEnvKey),
				port:     mustAtoi(t, mustGetEnv(t, mailerSMTPPortEnvKey)),
			},
		},
		{
			name:    "Should return error when failed to connect",
			wantErr: true,
			args: args{
				ctx: context.Background(),
				in: &pb.Mail{
					To:      []string{"test@gmail.com", "test2@gmail.com"},
					Html:    "hello",
					Subject: "world",
				},
			},
			fields: fields{
				from:     mustGetEnv(t, mailerSMTPFromEnvKey),
				password: mustGetEnv(t, mailerSMTPPasswordEnvKey),
			},
		},
	}

	mh := mailpit.NewClient(
		mustGetEnv(t, mailerSMTPHostEnvKey),
		mustAtoi(t, mustGetEnv(t, mailerTestAPIPortEnvKey)),
		time.Second*3,
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				require.NoError(t, mh.DeleteAll(), "Failed to cleanup emails")
			})

			c := NewClient(tt.fields.from, tt.fields.password, tt.fields.host, tt.fields.port)
			err := c.Send(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			msg, err := mh.GetAll()
			require.NoError(t, err)
			require.NotZero(t, len(msg))

			mail := msg[0]
			require.NoError(t, err)
			require.Empty(t, mail.Cc)
			require.Empty(t, mail.To)
			require.Equal(t, tt.args.in.GetTo(), getToMails(mail.Bcc))
			require.Contains(t, mail.Subject, tt.args.in.GetSubject())
		})
	}
}

func getToMails(to []mailpit.Receipient) []string {
	mails := make([]string, 0, len(to))
	for _, t := range to {
		mails = append(mails, t.Address)
	}
	return mails
}

func mustAtoi(t *testing.T, str string) int {
	t.Helper()
	num, err := strconv.Atoi(str)
	require.NoError(t, err)
	return num
}

func mustGetEnv(t *testing.T, key string) string {
	t.Helper()
	env := os.Getenv(key)
	require.NotEmpty(t, env)
	return env
}

func newImmediateCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	return ctx
}
