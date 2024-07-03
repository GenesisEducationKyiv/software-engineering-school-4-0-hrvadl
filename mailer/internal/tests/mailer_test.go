//go:build integration

package tests

import (
	"log/slog"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/pkg/mailpit"
	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v2/mailer"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/app"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/cfg"
)

const (
	mailerSMTPHostEnvKey     = "MAILER_TEST_SMTP_HOST"
	mailerSMTPPasswordEnvKey = "MAILER_TEST_SMTP_PASSWORD" // #nosec G101
	mailerSMTPFromEnvKey     = "MAILER_TEST_SMTP_FROM"
	mailerSMTPPortEnvKey     = "MAILER_TEST_SMTP_PORT"
	mailerTestAPIPortEnvKey  = "MAILER_TEST_API_PORT"
	mailerTestNatsURL        = "MAILER_NATS_TEST_URL"
)

// NOTE: I'm not doing component test using Resend mailer.
// It has quota and I don't wanna exceed it.
func TestAppRun(t *testing.T) {
	const topic = "mail"
	type args struct {
		subject string
		command *pb.MailCommand
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should send emails after receiving mail from topic",
			args: args{
				subject: topic,
				command: &pb.MailCommand{
					EventID:   "1",
					EventType: "sendMsg",
					Data: &pb.Mail{
						To:      []string{"test@test.com"},
						Subject: "should send to one receiver",
						Html:    "test HTML",
					},
				},
			},
		},
		{
			name: "Should send emails to multiple receivers after receiving mail from topic",
			args: args{
				subject: topic,
				command: &pb.MailCommand{
					EventID:   "1",
					EventType: "sendMsg",
					Data: &pb.Mail{
						To:      []string{"test@test.com", "test2@test.com"},
						Subject: "should send bcc to 2 receivers",
						Html:    "test HTML",
					},
				},
			},
		},
		{
			name: "Should return error when mail is not valid",
			args: args{
				subject: topic,
				command: &pb.MailCommand{
					EventID:   "1",
					EventType: "sendMsg",
					Data: &pb.Mail{
						To:      []string{"testest.com"},
						Subject: "should not send when it's not valid",
						Html:    "test HTML",
					},
				},
			},
			wantErr: true,
		},
	}

	cfg := cfg.Config{
		MailerFrom:     mustGetEnv(t, mailerSMTPFromEnvKey),
		MailerHost:     mustGetEnv(t, mailerSMTPHostEnvKey),
		MailerPassword: mustGetEnv(t, mailerSMTPPortEnvKey),
		NatsURL:        mustGetEnv(t, mailerTestNatsURL),
		MailerPort:     mustGetIntEnv(t, mailerSMTPPortEnvKey),
		Port:           "33300",
		Host:           "localhost",
	}

	mp := mailpit.NewClient(cfg.MailerHost, mustGetIntEnv(t, mailerTestAPIPortEnvKey), time.Second)
	nc, err := nats.Connect(cfg.NatsURL)
	require.NoError(t, err, "Failed to connect to NATS")

	app := app.New(cfg, slog.Default())
	go app.MustRun()
	time.Sleep(time.Second)
	t.Cleanup(func() {
		app.Stop()
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				require.NoError(t, mp.DeleteAll(), "Failed to cleanup emails")
			})

			bytes, err := proto.Marshal(tt.args.command)
			require.NoError(t, err)
			_, err = nc.Request(tt.args.subject, bytes, time.Second)
			require.NoError(t, err)

			time.Sleep(time.Second)
			data := tt.args.command.GetData()
			messages, err := mp.GetAll()
			require.NoError(t, err)
			if tt.wantErr {
				require.Len(t, messages, 0)
				return
			}

			require.Len(t, messages, 1)
			m := messages[0]
			require.Equal(t, data.GetSubject(), m.Subject)
			require.Equal(t, data.GetTo(), getToMails(m.Bcc))
		})
	}
}

func mustGetIntEnv(t *testing.T, key string) int {
	t.Helper()
	i, err := strconv.Atoi(mustGetEnv(t, key))
	require.NoError(t, err)
	return i
}

func mustGetEnv(t *testing.T, key string) string {
	t.Helper()
	env := os.Getenv(key)
	require.NotEmpty(t, env)
	return env
}

func getToMails(to []mailpit.Receipient) []string {
	mails := make([]string, 0, len(to))
	for _, t := range to {
		mails = append(mails, t.Address)
	}
	return mails
}
