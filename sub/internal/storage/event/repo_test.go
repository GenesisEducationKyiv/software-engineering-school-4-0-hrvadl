//go:build !integration

package event

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
)

func TestNewRepo(t *testing.T) {
	t.Parallel()
	type args struct {
		db *db.Tx
	}
	tests := []struct {
		name string
		args args
		want *Repo
	}{
		{
			name: "Should construct new repo with correct parameters",
			args: args{
				db: &db.Tx{},
			},
			want: &Repo{
				db: &db.Tx{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewRepo(tt.args.db)
			require.Equal(t, tt.want, got)
		})
	}
}
