package mailer

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/transport/queue/publisher/mailer/mocks"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	type args struct {
		pub Publisher
		log *slog.Logger
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "Should create new client correctly",
			args: args{
				pub: mocks.NewMockPublisher(gomock.NewController(t)),
				log: slog.Default(),
			},
			want: &Client{
				pub: mocks.NewMockPublisher(gomock.NewController(t)),
				log: slog.Default(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewClient(tt.args.pub, tt.args.log)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestClientSend(t *testing.T) {
	t.Parallel()
	type fields struct {
		log *slog.Logger
		pub func(c *gomock.Controller) Publisher
	}
	type args struct {
		ctx     context.Context
		html    string
		subject string
		to      []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should not return error when everything is correct",
			fields: fields{
				log: slog.Default(),
				pub: func(c *gomock.Controller) Publisher {
					pub := mocks.NewMockPublisher(c)
					pub.EXPECT().Publish(subject, gomock.Any()).Times(1).Return(nil)
					return pub
				},
			},
			args: args{
				ctx:     context.Background(),
				html:    "<h1>hello</h1>",
				subject: "sub",
				to:      []string{"to@to.com"},
			},
			wantErr: false,
		},
		{
			name: "Should return error when queue returned error",
			fields: fields{
				log: slog.Default(),
				pub: func(c *gomock.Controller) Publisher {
					pub := mocks.NewMockPublisher(c)
					pub.EXPECT().
						Publish(subject, gomock.Any()).
						Times(1).
						Return(errors.New("failed to send"))
					return pub
				},
			},
			args: args{
				ctx:     context.Background(),
				html:    "<h1>hello</h1>",
				subject: "sub",
				to:      []string{"to@to.com"},
			},
			wantErr: true,
		},
		{
			name: "Should return error when queue times out",
			fields: fields{
				log: slog.Default(),
				pub: func(c *gomock.Controller) Publisher {
					pub := mocks.NewMockPublisher(c)
					pub.EXPECT().
						Publish(subject, gomock.Any()).
						MaxTimes(1).
						Return(errors.New("failed to send"))
					return pub
				},
			},
			args: args{
				ctx:     newImmediateCtx(),
				html:    "<h1>hello</h1>",
				subject: "sub",
				to:      []string{"to@to.com"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Client{
				log: tt.fields.log,
				pub: tt.fields.pub(gomock.NewController(t)),
			}

			err := c.Send(tt.args.ctx, tt.args.html, tt.args.subject, tt.args.to...)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func newImmediateCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	return ctx
}
