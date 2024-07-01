//go:build !integration

package rateapi

import (
	"testing"

	"github.com/stretchr/testify/require"
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
