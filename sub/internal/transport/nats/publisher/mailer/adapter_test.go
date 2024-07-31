package mailer

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

func TestNewAdapter(t *testing.T) {
	t.Parallel()
	type args struct {
		nats    *nats.Conn
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want *Adapter
	}{
		{
			name: "Should construct adapter correctly",
			args: args{
				nats:    &nats.Conn{},
				timeout: time.Second * 3,
			},
			want: &Adapter{
				nats:    &nats.Conn{},
				timeout: time.Second * 3,
			},
		},
		{
			name: "Should construct adapter correctly",
			args: args{
				nats: &nats.Conn{
					Opts: nats.Options{
						Url: "nats://127.3.3.3:7767",
					},
				},
				timeout: time.Second * 3,
			},
			want: &Adapter{
				nats: &nats.Conn{
					Opts: nats.Options{
						Url: "nats://127.3.3.3:7767",
					},
				},
				timeout: time.Second * 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewAdapter(tt.args.nats, tt.args.timeout)
			require.Equal(t, tt.want, got)
		})
	}
}
