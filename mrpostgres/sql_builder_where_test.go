package mrpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the mrpostgres.SQLBuilderWhere conforms with the mrstorage.SQLBuilderWhere interface.
func TestSQLBuilderWhereImplementsSQLBuilderWhere(t *testing.T) {
	assert.Implements(t, (*mrstorage.SQLBuilderWhere)(nil), &mrpostgres.SQLBuilderWhere{})
}
