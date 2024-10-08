package mrpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the mrpostgres.SQLBuilderOrderBy conforms with the mrstorage.SQLBuilderOrderBy interface.
func TestSQLBuilderOrderByImplementsSQLBuilderOrderBy(t *testing.T) {
	assert.Implements(t, (*mrstorage.SQLBuilderOrderBy)(nil), &mrpostgres.SQLBuilderOrderBy{})
}
