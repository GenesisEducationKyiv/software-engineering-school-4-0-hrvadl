//go:build !integration

package subscriber

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/event"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/transaction"
)

func TestNewWithEventAdapter(t *testing.T) {
	t.Parallel()
	type args struct {
		r  *Repo
		er *event.Repo
		tx *transaction.Manager
	}
	tests := []struct {
		name string
		args args
		want *WithEventAdapter
	}{
		{
			name: "Should create new event adapter correctly",
			args: args{
				r:  &Repo{},
				er: &event.Repo{},
				tx: &transaction.Manager{},
			},
			want: &WithEventAdapter{
				repo:   &Repo{},
				events: &event.Repo{},
				tx:     &transaction.Manager{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewWithEventAdapter(tt.args.r, tt.args.er, tt.args.tx)
			require.Equal(t, tt.want, got)
		})
	}
}
