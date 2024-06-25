package sub

import (
	"context"
	"errors"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/sub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	subSvc "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/sub"
)

func NewErrorMappingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		res, err := handler(ctx, req)
		if errors.Is(err, subSvc.ErrInvalidEmail) {
			return res, invalidEmailToGRPCStatus(err)
		}

		if errors.Is(err, subSvc.ErrAlreadyExists) {
			return res, alreadyExistsToGRPCStatus(err)
		}

		return res, err
	}
}

func invalidEmailToGRPCStatus(err error) error {
	s := status.New(codes.InvalidArgument, err.Error())
	badReq := &pb.BadRequest{
		Code:        pb.ErrorCode_ERROR_CODE_INVALID_EMAIL_FORMAT,
		Field:       "email",
		Description: err.Error(),
	}

	s, err = s.WithDetails(badReq)
	if err != nil {
		return err
	}

	return s.Err()
}

func alreadyExistsToGRPCStatus(err error) error {
	s := status.New(codes.Aborted, err.Error())
	badReq := &pb.BadRequest{
		Code:        pb.ErrorCode_ERROR_CODE_ALREADY_EXISTS,
		Field:       "email",
		Description: err.Error(),
	}

	s, err = s.WithDetails(badReq)
	if err != nil {
		return err
	}

	return s.Err()
}
