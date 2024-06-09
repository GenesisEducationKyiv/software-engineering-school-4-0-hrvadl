//go:generate mockgen -destination=./mocks/mock_parser.go -package=mocks github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/mailer MailerServiceClient
package mailer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/mailer"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/pkg/logger"
)

const (
	retryCount   = 3
	retryTimeout = time.Second * 2
)

// NewClient constructs a GRPC mailer client with provided arguments. Under the hood
// it initializes a bunch of GRPC middleware for debugging and monitoring purposes. I.E:
// - retry middleware
// - request logger middleware
// If initialization of connection has failed it will return an error.
// NOTE: neither of parameters couldn't be nil or client will panic.
func NewClient(addr string, from string, log *slog.Logger) (*Client, error) {
	retryOpt := []retry.CallOption{
		retry.WithCodes(codes.Aborted, codes.NotFound, codes.DeadlineExceeded),
		retry.WithMax(retryCount),
		retry.WithPerRetryTimeout(retryTimeout),
	}

	cc, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			logger.NewClientGRPCMiddleware(log),
			retry.UnaryClientInterceptor(retryOpt...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sender service: %w", err)
	}

	return &Client{
		api:  pb.NewMailerServiceClient(cc),
		log:  log,
		from: from,
	}, nil
}

// Client represents GRPC mailer client which
// is responsible for sending messages to the provided
// subscribers. from parameter is author of the message, which is hard-coded
// on structure creation.
type Client struct {
	log  *slog.Logger
	api  pb.MailerServiceClient
	from string
}

func (c *Client) Send(ctx context.Context, html, subject string, to ...string) error {
	_, err := c.api.Send(ctx, &pb.Mail{
		From:    c.from,
		To:      to,
		Subject: subject,
		Html:    html,
	})
	return err
}
