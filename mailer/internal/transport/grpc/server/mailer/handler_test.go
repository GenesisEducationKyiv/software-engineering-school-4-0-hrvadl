//go:build !integration

package mailer

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/mailer"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/models/mail"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/transport/grpc/server/mailer/mocks"
)

func TestServerSend(t *testing.T) {
	t.Parallel()
	type fields struct {
		log    *slog.Logger
		client Client
	}
	type args struct {
		ctx context.Context
		m   *pb.Mail
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, client Client)
		want    *emptypb.Empty
		wantErr bool
	}{
		{
			name: "Should not return error when mail sender succeeded",
			fields: fields{
				log:    slog.Default(),
				client: mocks.NewMockClient(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				m: &pb.Mail{
					To:      []string{"to@to.com", "to1@to.com"},
					Html:    "test html",
					Subject: "test subject",
				},
			},
			setup: func(t *testing.T, client Client) {
				t.Helper()
				c, ok := client.(*mocks.MockClient)
				require.True(t, ok, "Failed to cast client to mock client")
				c.EXPECT().Send(gomock.Any(), mail.Mail{
					To:      []string{"to@to.com", "to1@to.com"},
					HTML:    "test html",
					Subject: "test subject",
				}).Times(1).Return(nil)
			},
			want:    &emptypb.Empty{},
			wantErr: false,
		},
		{
			name: "Should not return error when mail sender succeeded",
			fields: fields{
				log:    slog.Default(),
				client: mocks.NewMockClient(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				m: &pb.Mail{
					To:      []string{"to@to.com", "to1@to.com"},
					Html:    "test html",
					Subject: "test subject",
				},
			},
			setup: func(t *testing.T, client Client) {
				t.Helper()
				c, ok := client.(*mocks.MockClient)
				require.True(t, ok, "Failed to cast client to mock client")
				c.EXPECT().Send(gomock.Any(), mail.Mail{
					To:      []string{"to@to.com", "to1@to.com"},
					HTML:    "test html",
					Subject: "test subject",
				}).Times(1).Return(errors.New("failed to send mail"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup(t, tt.fields.client)
			s := &Server{
				UnimplementedMailerServiceServer: pb.UnimplementedMailerServiceServer{},
				log:                              tt.fields.log,
				client:                           tt.fields.client,
			}
			got, err := s.Send(tt.args.ctx, tt.args.m)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}
