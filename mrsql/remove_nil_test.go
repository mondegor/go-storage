package mrsql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

func TestSQLPartFuncRemoveNil(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		parts      []mrstorage.SQLPartFunc
		wantLength int
	}{
		{
			name: "test 1",
			parts: []mrstorage.SQLPartFunc{
				func(argumentNumber int) (sql string, args []any) { return "", nil },
				func(argumentNumber int) (sql string, args []any) { return "", nil },
			},
			wantLength: 2,
		},
		{
			name: "test 2",
			parts: []mrstorage.SQLPartFunc{
				func(argumentNumber int) (sql string, args []any) { return "", nil },
				nil,
				func(argumentNumber int) (sql string, args []any) { return "", nil },
			},
			wantLength: 2,
		},
		{
			name: "test 3",
			parts: []mrstorage.SQLPartFunc{
				nil,
				nil,
				nil,
			},
			wantLength: 0,
		},
		{
			name: "test 4",
			parts: []mrstorage.SQLPartFunc{
				func(argumentNumber int) (sql string, args []any) { return "", nil },
				nil,
				nil,
			},
			wantLength: 1,
		},
		{
			name: "test 5",
			parts: []mrstorage.SQLPartFunc{
				nil,
				nil,
				func(argumentNumber int) (sql string, args []any) { return "", nil },
				nil,
			},
			wantLength: 1,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := mrsql.SQLPartFuncRemoveNil(tt.parts)
			assert.Len(t, got, tt.wantLength)
		})
	}
}
