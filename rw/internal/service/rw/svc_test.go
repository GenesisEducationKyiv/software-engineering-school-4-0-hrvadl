package rw

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/rw/internal/service/rw/mocks"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	type args struct {
		source RateSource
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Should construct valid service",
			args: args{
				source: mocks.NewMockRateSource(gomock.NewController(t)),
			},
			want: &Service{
				sources: []RateSource{mocks.NewMockRateSource(gomock.NewController(t))},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewService(tt.args.source)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestServiceConvert(t *testing.T) {
	t.Parallel()
	type fields struct {
		sources func(ctrl *gomock.Controller) []RateSource
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float32
		wantErr bool
	}{
		{
			name: "Should not fallback when primary source succeeded",
			args: args{
				context.Background(),
			},
			fields: fields{
				sources: func(ctrl *gomock.Controller) []RateSource {
					rs1 := mocks.NewMockRateSource(ctrl)
					rs1.EXPECT().Convert(gomock.Any()).Times(1).Return(float32(3.3), nil)

					rs2 := mocks.NewMockRateSource(ctrl)
					rs2.EXPECT().Convert(gomock.Any()).Times(0)

					rs3 := mocks.NewMockRateSource(ctrl)
					rs3.EXPECT().Convert(gomock.Any()).Times(0)

					return []RateSource{rs1, rs2, rs3}
				},
			},
			want:    3.3,
			wantErr: false,
		},
		{
			name: "Should fallback to second when primary source failed",
			args: args{
				context.Background(),
			},
			fields: fields{
				sources: func(ctrl *gomock.Controller) []RateSource {
					rs1 := mocks.NewMockRateSource(ctrl)
					rs1.EXPECT().
						Convert(gomock.Any()).
						Times(1).
						Return(float32(0), errors.New("fail"))

					rs2 := mocks.NewMockRateSource(ctrl)
					rs2.EXPECT().Convert(gomock.Any()).Times(1).Return(float32(3.3), nil)

					rs3 := mocks.NewMockRateSource(ctrl)
					rs3.EXPECT().Convert(gomock.Any()).Times(0)

					return []RateSource{rs1, rs2, rs3}
				},
			},
			want:    3.3,
			wantErr: false,
		},
		{
			name: "Should fallback to third when primary & secondary source failed",
			args: args{
				context.Background(),
			},
			fields: fields{
				sources: func(ctrl *gomock.Controller) []RateSource {
					rs1 := mocks.NewMockRateSource(ctrl)
					rs1.EXPECT().
						Convert(gomock.Any()).
						Times(1).
						Return(float32(0), errors.New("fail"))

					rs2 := mocks.NewMockRateSource(ctrl)
					rs2.EXPECT().Convert(gomock.Any()).
						Times(1).
						Return(float32(0), errors.New("fail"))

					rs3 := mocks.NewMockRateSource(ctrl)
					rs3.EXPECT().Convert(gomock.Any()).Times(1).Return(float32(3.3), nil)

					return []RateSource{rs1, rs2, rs3}
				},
			},
			want:    3.3,
			wantErr: false,
		},
		{
			name: "Should return error when all services failed",
			args: args{
				context.Background(),
			},
			fields: fields{
				sources: func(ctrl *gomock.Controller) []RateSource {
					rs1 := mocks.NewMockRateSource(ctrl)
					rs1.EXPECT().
						Convert(gomock.Any()).
						Times(1).
						Return(float32(0), errors.New("fail"))

					rs2 := mocks.NewMockRateSource(ctrl)
					rs2.EXPECT().Convert(gomock.Any()).
						Times(1).
						Return(float32(0), errors.New("fail"))

					rs3 := mocks.NewMockRateSource(ctrl)
					rs3.EXPECT().Convert(gomock.Any()).
						Times(1).
						Return(float32(0), errors.New("fail"))

					return []RateSource{rs1, rs2, rs3}
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &Service{
				sources: tt.fields.sources(gomock.NewController(t)),
			}

			got, err := s.Convert(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.InDelta(t, tt.want, got, 2)
		})
	}
}

func TestServiceSetNext(t *testing.T) {
	t.Parallel()
	type fields struct {
		sources []RateSource
	}
	type args struct {
		source []RateSource
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Service
	}{
		{
			name: "Should set next service correctly",
			fields: fields{
				sources: []RateSource{
					mocks.NewMockRateSource(gomock.NewController(t)),
				},
			},
			args: args{
				source: []RateSource{mocks.NewMockRateSource(gomock.NewController(t))},
			},
			want: &Service{
				sources: []RateSource{
					mocks.NewMockRateSource(gomock.NewController(t)),
					mocks.NewMockRateSource(gomock.NewController(t)),
				},
			},
		},
		{
			name: "Should set next service correctly",
			fields: fields{
				sources: []RateSource{
					mocks.NewMockRateSource(gomock.NewController(t)),
				},
			},
			args: args{
				source: []RateSource{
					mocks.NewMockRateSource(gomock.NewController(t)),
					mocks.NewMockRateSource(gomock.NewController(t)),
				},
			},
			want: &Service{
				sources: []RateSource{
					mocks.NewMockRateSource(gomock.NewController(t)),
					mocks.NewMockRateSource(gomock.NewController(t)),
					mocks.NewMockRateSource(gomock.NewController(t)),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := &Service{
				sources: tt.fields.sources,
			}
			got.SetNext(tt.args.source...)
			require.Equal(t, tt.want, got)
		})
	}
}
