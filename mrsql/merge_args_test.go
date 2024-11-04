package mrsql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrsql"
)

func TestMergeArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args [][]any
		want []any
	}{
		{
			name: "test 1",
			args: [][]any{
				{1, 2, 3},
				{3, 4, 5},
			},
			want: []any{1, 2, 3, 3, 4, 5},
		},
		{
			name: "test 2",
			args: nil,
			want: []any{},
		},
		{
			name: "test 3",
			args: [][]any{
				{},
			},
			want: []any{},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := mrsql.MergeArgs(tt.args...)
			assert.Equal(t, tt.want, got)
		})
	}
}
