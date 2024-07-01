//go:build !integration

package privat

import (
	"testing"

	"github.com/stretchr/testify/require"
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
