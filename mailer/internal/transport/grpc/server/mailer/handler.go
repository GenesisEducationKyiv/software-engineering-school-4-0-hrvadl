package mailer

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/mailer"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/mail"
)

const operation = "mailer server"

// Register registers subscribe handler to the given GRPC server.
// NOTE: all parameters are required, the service will panic if
// either of them is missing.
func Register(srv *grpc.Server, client Client, log *slog.Logger) {
	pb.RegisterMailerServiceServer(srv, &Server{
		log:    log,
		client: client,
	})
}

//go:generate mockgen -destination=./mocks/mock_client.go -package=mocks . Client
type Client interface {
	Send(ctx context.Context, m mail.Mail) error
}

// Server represents mailer GRPC server
// which will handle the incoming requests and delegate
// all work to the underlying client.
type Server struct {
	pb.UnimplementedMailerServiceServer
	log    *slog.Logger
	client Client
}

// Send method calls underlying client method and returns an error, in case there was a
// failure.
func (s *Server) Send(ctx context.Context, in *pb.Mail) (*emptypb.Empty, error) {
	m := mail.Mail{
		To:      in.GetTo(),
		Subject: in.GetSubject(),
		HTML:    in.GetHtml(),
	}
	if err := s.client.Send(ctx, m); err != nil {
		return nil, fmt.Errorf("%s: failed to send mail: %w", operation, err)
	}
	return &emptypb.Empty{}, nil
}
