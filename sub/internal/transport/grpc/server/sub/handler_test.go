package sub

import (
	"errors"
	"log/slog"
	"reflect"
	"testing"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/sub"
	"go.uber.org/mock/gomock"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/transport/grpc/server/sub/mocks"
)

func TestServerSubscribe(t *testing.T) {
	t.Parallel()
	type fields struct {
		log *slog.Logger
		svc Service
	}
	type args struct {
		ctx context.Context
		req *pb.SubscribeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, svc Service)
		want    *emptypb.Empty
		wantErr bool
	}{
		{
			name: "Should not return any error when service succeded",
			fields: fields{
				log: slog.Default(),
				svc: mocks.NewMockService(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				req: &pb.SubscribeRequest{Email: "test@test.com"},
			},
			setup: func(t *testing.T, svc Service) {
				t.Helper()
				s, ok := svc.(*mocks.MockService)
				if !ok {
					t.Fatalf("Failed to cast service to mock service")
				}

				s.EXPECT().Subscribe(gomock.Any(), "test@test.com").Times(1).Return(int64(1), nil)
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Should return any error when service failed",
			fields: fields{
				log: slog.Default(),
				svc: mocks.NewMockService(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				req: &pb.SubscribeRequest{Email: "test@test.com"},
			},
			setup: func(t *testing.T, svc Service) {
				t.Helper()
				s, ok := svc.(*mocks.MockService)
				if !ok {
					t.Fatalf("Failed to cast service to mock service")
				}

				s.EXPECT().
					Subscribe(gomock.Any(), "test@test.com").
					Times(1).
					Return(int64(0), errors.New("failed to subscribe"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup(t, tt.fields.svc)
			s := &Server{
				UnimplementedSubServiceServer: pb.UnimplementedSubServiceServer{},
				log:                           tt.fields.log,
				svc:                           tt.fields.svc,
			}
			got, err := s.Subscribe(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.Subscribe() = %v, want %v", got, tt.want)
			}
		})
	}
}
