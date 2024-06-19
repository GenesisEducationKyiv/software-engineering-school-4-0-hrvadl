//go:build !integration

package gomail

import (
	"context"
	"errors"
	"reflect"
	"testing"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/mailer"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gopkg.in/gomail.v2"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/platform/mail/resend/mocks"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	type args struct {
		from     string
		password string
		host     string
		port     int
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "Should create client correctly",
			args: args{
				from:     "from222@again.com",
				password: "test",
				host:     "host3.com",
				port:     666,
			},
			want: &Client{
				dialer: gomail.NewDialer("host3.com", 666, "from222@again.com", "test"),
				from:   "from222@again.com",
			},
		},
		{
			name: "Should create client correctly",
			args: args{
				from:     "from@from.com",
				password: "test",
				host:     "host.com",
				port:     444,
			},
			want: &Client{
				dialer: gomail.NewDialer("host.com", 444, "from@from.com", "test"),
				from:   "from@from.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewClient(tt.args.from, tt.args.password, tt.args.host, tt.args.port); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientSend(t *testing.T) {
	t.Parallel()
	type fields struct {
		dialer *gomail.Dialer
		from   string
		next   ChainedSender
	}
	type args struct {
		ctx context.Context
		in  *pb.Mail
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, sender ChainedSender)
		wantErr bool
	}{
		{
			name: "Should fallback when failed",
			args: args{
				ctx: context.Background(),
				in: &pb.Mail{
					To:      []string{"to@to.com", "to1@to.com"},
					Html:    "html",
					Subject: "subject",
				},
			},
			fields: fields{
				dialer: gomail.NewDialer("test.com", 222, "from@from.com", "secret"),
				from:   "from@from.com",
				next:   mocks.NewMockChainedSender(gomock.NewController(t)),
			},
			wantErr: false,
			setup: func(t *testing.T, sender ChainedSender) {
				t.Helper()
				s, ok := sender.(*mocks.MockChainedSender)
				require.True(t, ok, "Failed to cast sender to mock")
				s.EXPECT().Send(gomock.Any(), &pb.Mail{
					To:      []string{"to@to.com", "to1@to.com"},
					Html:    "html",
					Subject: "subject",
				}).Times(1).Return(nil)
			},
		},
		{
			name: "Should fallback when failed",
			args: args{
				ctx: context.Background(),
				in: &pb.Mail{
					To:      []string{"to@to.com", "to1@to.com"},
					Html:    "html",
					Subject: "subject",
				},
			},
			fields: fields{
				dialer: gomail.NewDialer("test.com", 222, "from@from.com", "secret"),
				from:   "from@from.com",
				next:   mocks.NewMockChainedSender(gomock.NewController(t)),
			},
			wantErr: true,
			setup: func(t *testing.T, sender ChainedSender) {
				t.Helper()
				s, ok := sender.(*mocks.MockChainedSender)
				require.True(t, ok, "Failed to cast sender to mock")
				s.EXPECT().Send(gomock.Any(), &pb.Mail{
					To:      []string{"to@to.com", "to1@to.com"},
					Html:    "html",
					Subject: "subject",
				}).Times(1).Return(errors.New("failed to send"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup(t, tt.fields.next)
			c := &Client{
				dialer: tt.fields.dialer,
				from:   tt.fields.from,
				next:   tt.fields.next,
			}

			err := c.Send(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestClientSetNext(t *testing.T) {
	t.Parallel()
	type fields struct {
		dialer *gomail.Dialer
		from   string
		next   ChainedSender
	}
	type args struct {
		next ChainedSender
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Should set next correctly",
			args: args{
				next: mocks.NewMockChainedSender(gomock.NewController(t)),
			},
			fields: fields{
				dialer: gomail.NewDialer("test.com", 222, "test", "test"),
				from:   "test.com",
				next:   nil,
			},
		},
		{
			name: "Should set next correctly",
			args: args{
				next: nil,
			},
			fields: fields{
				dialer: gomail.NewDialer("test.com", 222, "test", "test"),
				from:   "test.com",
				next:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Client{
				dialer: tt.fields.dialer,
				from:   tt.fields.from,
				next:   tt.fields.next,
			}
			c.SetNext(tt.args.next)
			require.Equal(t, tt.args.next, c.next)
		})
	}
}
