//go:build integration

package exchangeapi

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/rw/internal/platform/rates/privat"
)

const (
	exchangeTestAPIBaseURLEnvKey         = "EXCHANGE_TEST_API_BASE_URL"
	exchangeTestAPIFallbackBaseURLEnvKey = "EXCHANGE_TEST_API_FALLBACK_BASE_URL"
)

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
		wantErr bool
	}{
		{
			name: "Should fallback to second exchanger",
			fields: fields{
				next: privat.NewClient(mustGetEnv(t, exchangeTestAPIFallbackBaseURLEnvKey)),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name:   "Should not get exchange rate correctly when URL is missing",
			fields: fields{},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "Should not get exchange rate correctly when context exceeded",
			fields: fields{
				url: mustGetEnv(t, exchangeTestAPIBaseURLEnvKey),
			},
			args: args{
				ctx: newImmediateCtx(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := Client{
				url:  tt.fields.url,
				next: tt.fields.next,
			}

			got, err := c.Convert(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotZero(t, got)
		})
	}
}

func mustGetEnv(t *testing.T, key string) string {
	t.Helper()
	env := os.Getenv(key)
	require.NotEmpty(t, env)
	return env
}

func newImmediateCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	return ctx
}
