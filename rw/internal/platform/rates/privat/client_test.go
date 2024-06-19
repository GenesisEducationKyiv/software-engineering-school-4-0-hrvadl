//go:build !integration

package privat

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/rw/internal/platform/rates/privat/mocks"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "Should construct client correctly",
			args: args{
				url: "https://privat.bank.com",
			},
			want: &Client{
				url: "https://privat.bank.com",
			},
		},
		{
			name: "Should construct client correctly",
			args: args{
				url: "https://privat.bank1.com",
			},
			want: &Client{
				url: "https://privat.bank1.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewClient(tt.args.url)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestClientConvert(t *testing.T) {
	t.Parallel()
	type fields struct {
		url  string
		next Converter
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
				url:  "https://url.com",
				next: mocks.NewMockConverter(gomock.NewController(t)),
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
				url:  "https://exchange.com",
				next: mocks.NewMockConverter(gomock.NewController(t)),
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
				url:  tt.fields.url,
				next: tt.fields.next,
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
