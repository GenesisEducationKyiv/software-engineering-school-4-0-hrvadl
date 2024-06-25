//go:build !integration

package rateapi

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/rw/internal/platform/rates/rateapi/mocks"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	type args struct {
		token string
		url   string
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "Should create new client correctly",
			args: args{
				token: "token",
				url:   "https://url.com",
			},
			want: &Client{
				token: "token",
				url:   "https://url.com",
			},
		},
		{
			name: "Should create new client correctly",
			args: args{
				token: "token2266",
				url:   "https://url2.com",
			},
			want: &Client{
				token: "token2266",
				url:   "https://url2.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, NewClient(tt.args.token, tt.args.url))
		})
	}
}

func TestClientConvert(t *testing.T) {
	t.Parallel()
	type fields struct {
		token string
		url   string
		next  Converter
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, converter Converter)
		want    float32
		wantErr bool
	}{
		{
			name: "Should fallback to another client when failed",
			fields: fields{
				token: "token",
				url:   "https://url.com",
				next:  mocks.NewMockConverter(gomock.NewController(t)),
			},
			setup: func(t *testing.T, converter Converter) {
				t.Helper()
				c, ok := converter.(*mocks.MockConverter)
				require.True(t, ok, "Failed to cast converter to mock")
				c.EXPECT().Convert(gomock.Any()).Times(1).Return(float32(22), nil)
			},
			want: 22,
		},
		{
			name: "Should fallback to another client when failed",
			fields: fields{
				token: "token",
				url:   "https://exchange.com",
				next:  mocks.NewMockConverter(gomock.NewController(t)),
			},
			args: args{
				ctx: newImmediateCtx(),
			},
			setup: func(t *testing.T, converter Converter) {
				t.Helper()
				c, ok := converter.(*mocks.MockConverter)
				require.True(t, ok, "Failed to cast converter to mock")
				c.EXPECT().
					Convert(gomock.Any()).
					Times(1).
					Return(float32(0), errors.New("failed to convert"))
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup(t, tt.fields.next)
			c := Client{
				url:   tt.fields.url,
				next:  tt.fields.next,
				token: tt.fields.token,
			}

			got, err := c.Convert(tt.args.ctx)
			require.InDelta(t, tt.want, got, 2)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
		})
	}
}

func newImmediateCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	return ctx
}
