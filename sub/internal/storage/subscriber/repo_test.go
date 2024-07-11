//go:build !integration

package subscriber

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRepo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
	}{
		{
			name: "Should create repo with correct db conn",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewRepo()
			require.NotNil(t, got)
		})
	}
}
