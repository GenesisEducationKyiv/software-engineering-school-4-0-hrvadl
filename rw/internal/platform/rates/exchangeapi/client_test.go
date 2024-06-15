package exchangeapi

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
		want Client
	}{
		{
			name: "Should initiate new client correctly",
			args: args{
				url: "https://url.com",
			},
			want: Client{
				url: "https://url.com",
			},
		},
		{
			name: "Should initiate new client correctly",
			args: args{
				url: "https://url2.com",
			},
			want: Client{
				url: "https://url2.com",
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
