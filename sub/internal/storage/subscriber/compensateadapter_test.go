package subscriber

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/transaction"
)

func TestNewCompensateAdapter(t *testing.T) {
	t.Parallel()
	type args struct {
		r  *Repo
		tx *transaction.Manager
	}
	tests := []struct {
		name string
		args args
		want *CompensateAdapter
	}{
		{
			name: "Should create compensate adapter correctly",
			args: args{
				r:  &Repo{},
				tx: &transaction.Manager{},
			},
			want: &CompensateAdapter{
				repo: &Repo{},
				tx:   &transaction.Manager{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewCompensateAdapter(tt.args.r, tt.args.tx)
			require.Equal(t, tt.want, got)
		})
	}
}
