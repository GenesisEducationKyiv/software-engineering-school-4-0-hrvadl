//go:build integration

package privat

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/rw/internal/platform/rates/exchangeapi"
)

const exchangeTestAPIBaseURLEnvKey = "EXCHANGE_TEST_API_BASE_URL"

func TestClientConvert(t *testing.T) {
	t.Parallel()
	type fields struct {
		next Converter
		url  string
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
			name: "Should convert correctly",
			fields: fields{
				url: "https://api.privatbank.ua",
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Should return error when request takes too long",
			fields: fields{
				url: "https://api.privatbank.ua",
			},
			args: args{
				ctx: newImmediateCtx(),
			},
			wantErr: true,
		},
		{
			name: "Should return error when url is incorrect",
			fields: fields{
				url: "https://api",
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "Should fallback to second exchanger",
			fields: fields{
				url:  "https://api.privatbank.ua",
				next: exchangeapi.NewClient(mustGetEnv(t, exchangeTestAPIBaseURLEnvKey)),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := Client{
				url: tt.fields.url,
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

func newImmediateCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	return ctx
}

func mustGetEnv(t *testing.T, key string) string {
	t.Helper()
	env := os.Getenv(key)
	require.NotEmpty(t, env)
	return env
}
