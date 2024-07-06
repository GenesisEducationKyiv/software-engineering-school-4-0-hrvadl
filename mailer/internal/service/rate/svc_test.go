//go:build !integration

package rate

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/service/rate/mocks"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/rate"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	type args struct {
		rs RateSource
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Should create new service correctly",
			args: args{
				rs: mocks.NewMockRateSource(gomock.NewController(t)),
			},
			want: &Service{
				rate: mocks.NewMockRateSource(gomock.NewController(t)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewService(tt.args.rs)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestServiceGet(t *testing.T) {
	t.Parallel()
	type fields struct {
		rate func(c *gomock.Controller) RateSource
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *rate.Exchange
		wantErr bool
	}{
		{
			name: "Should not return error when rate source succeeded",
			fields: fields{
				rate: func(c *gomock.Controller) RateSource {
					rs := mocks.NewMockRateSource(c)
					rs.EXPECT().Get(context.Background()).Return(&rate.Exchange{
						From: "USD",
						To:   "UAH",
						Rate: 33.3,
					}, nil)
					return rs
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
			want: &rate.Exchange{
				From: "USD",
				To:   "UAH",
				Rate: 33.3,
			},
		},
		{
			name: "Should not return error when rate source succeeded",
			fields: fields{
				rate: func(c *gomock.Controller) RateSource {
					rs := mocks.NewMockRateSource(c)
					rs.EXPECT().
						Get(context.Background()).
						Return(nil, errors.New("failed to get rate"))
					return rs
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Service{
				rate: tt.fields.rate(gomock.NewController(t)),
			}
			got, err := s.Get(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestServiceReplace(t *testing.T) {
	t.Parallel()
	type fields struct {
		rate func(c *gomock.Controller) RateSource
	}
	type args struct {
		ctx  context.Context
		rate rate.Exchange
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should not return error if rate source succeeded",
			fields: fields{
				rate: func(c *gomock.Controller) RateSource {
					rs := mocks.NewMockRateSource(c)
					rs.EXPECT().Replace(context.Background(), rate.Exchange{
						From: "USD",
						To:   "UAH",
						Rate: 44.4,
					}).Times(1).Return(nil)
					return rs
				},
			},
			args: args{
				ctx: context.Background(),
				rate: rate.Exchange{
					From: "USD",
					To:   "UAH",
					Rate: 44.4,
				},
			},
			wantErr: false,
		},
		{
			name: "Should return error if rate source failed",
			fields: fields{
				rate: func(c *gomock.Controller) RateSource {
					rs := mocks.NewMockRateSource(c)
					rs.EXPECT().Replace(context.Background(), rate.Exchange{
						From: "USD",
						To:   "UAH",
						Rate: 44.4,
					}).Times(1).Return(errors.New("failed to replace"))
					return rs
				},
			},
			args: args{
				ctx: context.Background(),
				rate: rate.Exchange{
					From: "USD",
					To:   "UAH",
					Rate: 44.4,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Service{
				rate: tt.fields.rate(gomock.NewController(t)),
			}
			err := s.Replace(tt.args.ctx, tt.args.rate)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
