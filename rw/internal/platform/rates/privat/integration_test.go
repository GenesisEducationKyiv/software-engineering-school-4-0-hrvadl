//go:build integration

package privat

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

const exchangeTestAPIBaseURLEnvKey = "EXCHANGE_TEST_API_BASE_URL"

func TestClientConvert(t *testing.T) {
	t.Parallel()
	type fields struct {
		url string
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
