package rates

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/rw/internal/platform/rates/mocks"
)

func TestNewWithLogger(t *testing.T) {
	t.Parallel()
	type args struct {
		base Converter
		log  *slog.Logger
	}
	tests := []struct {
		name string
		args args
		want *WithLoggerDecorator
	}{
		{
			name: "Should construct decorator correctly",
			args: args{
				base: mocks.NewMockConverter(gomock.NewController(t)),
				log:  slog.New(slog.NewJSONHandler(os.Stdout, nil)),
			},
			want: &WithLoggerDecorator{
				base: mocks.NewMockConverter(gomock.NewController(t)),
				log:  slog.New(slog.NewJSONHandler(os.Stdout, nil)),
			},
		},
		{
			name: "Should construct decorator correctly",
			args: args{
				base: mocks.NewMockConverter(gomock.NewController(t)),
				log:  slog.New(slog.NewJSONHandler(os.Stderr, nil)),
			},
			want: &WithLoggerDecorator{
				base: mocks.NewMockConverter(gomock.NewController(t)),
				log:  slog.New(slog.NewJSONHandler(os.Stderr, nil)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewWithLogger(tt.args.base, tt.args.log)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestWithLoggerDecoratorConvert(t *testing.T) {
	t.Parallel()
	type fields struct {
		base Converter
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		buf     *bytes.Buffer
		args    args
		want    float32
		setup   func(t *testing.T, converter Converter)
		wantErr bool
	}{
		{
			name: "Should log successful response",
			args: args{
				ctx: context.Background(),
			},

			fields: fields{
				base: mocks.NewMockConverter(gomock.NewController(t)),
			},
			buf: bytes.NewBuffer([]byte{}),
			setup: func(t *testing.T, converter Converter) {
				t.Helper()
				c, ok := converter.(*mocks.MockConverter)
				require.True(t, ok, "Failed to cast converter to mock")
				c.EXPECT().Convert(gomock.Any()).Times(1).Return(float32(3), nil)
			},
			want: 3,
		},
		{
			name: "Should log failed response",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				base: mocks.NewMockConverter(gomock.NewController(t)),
			},
			buf: bytes.NewBuffer([]byte{}),
			setup: func(t *testing.T, converter Converter) {
				t.Helper()
				c, ok := converter.(*mocks.MockConverter)
				require.True(t, ok, "Failed to cast converter to mock")
				c.EXPECT().
					Convert(gomock.Any()).
					Times(1).
					Return(float32(0), errors.New("failed to convert"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup(t, tt.fields.base)
			c := WithLoggerDecorator{
				base: tt.fields.base,
				log:  slog.New(slog.NewJSONHandler(tt.buf, nil)),
			}

			got, err := c.Convert(tt.args.ctx)
			require.InDelta(t, tt.want, got, 2)

			require.Contains(t, tt.buf.String(), "Sending request to exchange API service")
			if !tt.wantErr {
				require.Contains(t, tt.buf.String(), "Received response from exchange API")
				return
			}
			require.Error(t, err)
			require.Contains(t, tt.buf.String(), "Received error from exchange API")
		})
	}
}
