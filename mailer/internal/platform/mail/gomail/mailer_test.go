//go:build !integration

package gomail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/gomail.v2"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	type args struct {
		from     string
		password string
		host     string
		port     int
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "Should create client correctly",
			args: args{
				from:     "from222@again.com",
				password: "test",
				host:     "host3.com",
				port:     666,
			},
			want: &Client{
				dialer: gomail.NewDialer("host3.com", 666, "from222@again.com", "test"),
				from:   "from222@again.com",
			},
		},
		{
			name: "Should create client correctly",
			args: args{
				from:     "from@from.com",
				password: "test",
				host:     "host.com",
				port:     444,
			},
			want: &Client{
				dialer: gomail.NewDialer("host.com", 444, "from@from.com", "test"),
				from:   "from@from.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewClient(tt.args.from, tt.args.password, tt.args.host, tt.args.port)
			require.Equal(t, tt.want, got)
		})
	}
}
