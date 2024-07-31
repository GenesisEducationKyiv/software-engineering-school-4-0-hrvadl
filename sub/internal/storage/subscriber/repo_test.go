//go:build !integration

package subscriber

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
)

func TestNewRepo(t *testing.T) {
	type args struct {
		db *db.Tx
	}
	t.Parallel()
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should create repo with correct db conn",
			args: args{
				db: &db.Tx{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewRepo(tt.args.db)
			require.NotNil(t, got)
		})
	}
}
