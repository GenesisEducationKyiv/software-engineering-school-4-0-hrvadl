//go:build !integration

package resend

import (
	"context"
	"errors"
	"testing"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/mailer"
	rs "github.com/resend/resend-go/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/platform/mail/resend/mocks"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	type args struct {
		from  string
		token string
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "Should construct new client correctly",
			args: args{
				from:  "from@gmail.com",
				token: "secret",
			},
			want: &Client{
				from:   "from@gmail.com",
				client: rs.NewClient("secret"),
			},
		},
		{
			name: "Should construct new client correctly",
			args: args{
				from:  "from31313@gmail.com",
				token: "secret2222",
			},
			want: &Client{
				from:   "from31313@gmail.com",
				client: rs.NewClient("secret2222"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewClient(tt.args.from, tt.args.token)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestClientSend(t *testing.T) {
	t.Parallel()
	type fields struct {
		client *rs.Client
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
				client: rs.NewClient("tset"),
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
				client: rs.NewClient("tset"),
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
				from:   tt.fields.from,
				client: tt.fields.client,
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
		client *rs.Client
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
			fields: fields{
				client: rs.NewClient("key"),
				from:   "from",
				next:   nil,
			},
			args: args{
				next: mocks.NewMockChainedSender(gomock.NewController(t)),
			},
		},
		{
			name: "Should set next correctly",
			fields: fields{
				client: rs.NewClient("key"),
				from:   "from",
				next:   nil,
			},
			args: args{
				next: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Client{
				client: tt.fields.client,
				from:   tt.fields.from,
				next:   tt.fields.next,
			}
			c.SetNext(tt.args.next)
			require.Equal(t, tt.args.next, c.next)
		})
	}
}
