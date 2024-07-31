//go:build !integration

package subscriber

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewRepo(t *testing.T) {
	t.Parallel()
	type args struct {
		db *mongo.Database
	}
	tests := []struct {
		name string
		args args
		want *Repository
	}{
		{
			name: "Should construct new repo correctly",
			args: args{
				db: &mongo.Database{},
			},
			want: &Repository{
				db: &mongo.Database{},
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
