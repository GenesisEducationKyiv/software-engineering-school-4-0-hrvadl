//go:build !integration

package resend

import (
	"testing"

	rs "github.com/resend/resend-go/v2"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	type args struct {
		from  string
		token string
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "Should construct new client correctly",
			args: args{
				from:  "from@gmail.com",
				token: "secret",
			},
			want: &Client{
				from:   "from@gmail.com",
				client: rs.NewClient("secret"),
			},
		},
		{
			name: "Should construct new client correctly",
			args: args{
				from:  "from31313@gmail.com",
				token: "secret2222",
			},
			want: &Client{
				from:   "from31313@gmail.com",
				client: rs.NewClient("secret2222"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewClient(tt.args.from, tt.args.token)
			require.Equal(t, tt.want, got)
		})
	}
}
