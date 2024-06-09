//go:generate mockgen -destination=./mocks/mock_rw.go -package=mocks github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/ratewatcher RateWatcherServiceClient
package ratewatcher

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/ratewatcher"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/pkg/logger"
)

const (
	retryCount   = 3
	retryTimeout = time.Second * 2
)

// NewClient constructs a GRPC rate watcher client with provided arguments. Under the hood
// it initializes a bunch of GRPC middleware for debugging and monitoring purposes. I.E:
// - retry middleware
// - request logger middleware
// If initialization of connection has failed it will return an error.
// NOTE: neither of parameters couldn't be nil or client will panic.
func NewClient(addr string, log *slog.Logger) (*Client, error) {
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
		api: pb.NewRateWatcherServiceClient(cc),
	}, nil
}

// Client represents GRPC rate wathcer client which
// is responsible for getting latest exchange rates and
// returning it in response.
type Client struct {
	api pb.RateWatcherServiceClient
}

func (c *Client) GetRate(ctx context.Context) (float32, error) {
	resp, err := c.api.GetRate(ctx, &emptypb.Empty{})
	if err != nil {
		return 0, err
	}

	return resp.Rate, nil
}
