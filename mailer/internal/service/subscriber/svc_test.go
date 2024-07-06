package subscriber

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/service/subscriber/mocks"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/subscriber"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	type args struct {
		ss SubscriberSource
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Should create new service correctly",
			args: args{
				ss: mocks.NewMockSubscriberSource(gomock.NewController(t)),
			},
			want: &Service{
				subscriber: mocks.NewMockSubscriberSource(gomock.NewController(t)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewService(tt.args.ss)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestServiceGetAll(t *testing.T) {
	t.Parallel()
	type fields struct {
		subscriber func(c *gomock.Controller) SubscriberSource
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []subscriber.Subscriber
		wantErr bool
	}{
		{
			name: "Should not return error when source succeeded",
			fields: fields{
				subscriber: func(c *gomock.Controller) SubscriberSource {
					s := mocks.NewMockSubscriberSource(c)
					s.EXPECT().GetAll(context.Background()).Times(1).Return([]subscriber.Subscriber{
						{Email: "test@test.com"},
					}, nil)
					return s
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: []subscriber.Subscriber{
				{Email: "test@test.com"},
			},
			wantErr: false,
		},
		{
			name: "Should not return error when source succeeded",
			fields: fields{
				subscriber: func(c *gomock.Controller) SubscriberSource {
					s := mocks.NewMockSubscriberSource(c)
					s.EXPECT().GetAll(context.Background()).Times(1).Return([]subscriber.Subscriber{
						{Email: "test@test.com"},
						{Email: "test2@test.com"},
						{Email: "test3@test.com"},
					}, nil)
					return s
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: []subscriber.Subscriber{
				{Email: "test@test.com"},
				{Email: "test2@test.com"},
				{Email: "test3@test.com"},
			},
			wantErr: false,
		},
		{
			name: "Should return error when source failed",
			fields: fields{
				subscriber: func(c *gomock.Controller) SubscriberSource {
					s := mocks.NewMockSubscriberSource(c)
					s.EXPECT().GetAll(context.Background()).Times(1).Return(nil, ErrNoSubscribers)
					return s
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Service{
				subscriber: tt.fields.subscriber(gomock.NewController(t)),
			}
			got, err := s.GetAll(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestServiceSave(t *testing.T) {
	t.Parallel()
	type fields struct {
		subscriber func(c *gomock.Controller) SubscriberSource
	}
	type args struct {
		ctx context.Context
		sub subscriber.Subscriber
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should not return error when source succeeded",
			fields: fields{
				subscriber: func(c *gomock.Controller) SubscriberSource {
					ss := mocks.NewMockSubscriberSource(c)
					ss.EXPECT().
						Save(context.Background(), subscriber.Subscriber{Email: "test@test.com"}).
						Times(1).
						Return(nil)
					return ss
				},
			},
			args: args{
				ctx: context.Background(),
				sub: subscriber.Subscriber{
					Email: "test@test.com",
				},
			},
			wantErr: false,
		},
		{
			name: "Should return error when source failed",
			fields: fields{
				subscriber: func(c *gomock.Controller) SubscriberSource {
					ss := mocks.NewMockSubscriberSource(c)
					ss.EXPECT().
						Save(context.Background(), subscriber.Subscriber{Email: "test@test.com"}).
						Times(1).
						Return(errors.New("failed to save"))
					return ss
				},
			},
			args: args{
				ctx: context.Background(),
				sub: subscriber.Subscriber{
					Email: "test@test.com",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Service{
				subscriber: tt.fields.subscriber(gomock.NewController(t)),
			}
			err := s.Save(tt.args.ctx, tt.args.sub)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestServiceDelete(t *testing.T) {
	t.Parallel()
	type fields struct {
		subscriber func(c *gomock.Controller) SubscriberSource
	}
	type args struct {
		ctx context.Context
		sub subscriber.Subscriber
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should not return error when source succeeded",
			fields: fields{
				subscriber: func(c *gomock.Controller) SubscriberSource {
					ss := mocks.NewMockSubscriberSource(c)
					ss.EXPECT().
						Delete(context.Background(), subscriber.Subscriber{Email: "test@test.com"}).
						Times(1).
						Return(nil)
					return ss
				},
			},
			args: args{
				ctx: context.Background(),
				sub: subscriber.Subscriber{
					Email: "test@test.com",
				},
			},
			wantErr: false,
		},
		{
			name: "Should return error when source failed",
			fields: fields{
				subscriber: func(c *gomock.Controller) SubscriberSource {
					ss := mocks.NewMockSubscriberSource(c)
					ss.EXPECT().
						Delete(context.Background(), subscriber.Subscriber{Email: "test@test.com"}).
						Times(1).
						Return(errors.New("failed to delete"))
					return ss
				},
			},
			args: args{
				ctx: context.Background(),
				sub: subscriber.Subscriber{
					Email: "test@test.com",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Service{
				subscriber: tt.fields.subscriber(gomock.NewController(t)),
			}
			err := s.Delete(tt.args.ctx, tt.args.sub)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
