package mrpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the mrpostgres.SQLBuilderUpdate conforms with the mrstorage.SQLBuilderUpdate interface.
func TestSQLBuilderUpdateImplementsSQLBuilderUpdate(t *testing.T) {
	assert.Implements(t, (*mrstorage.SQLBuilderUpdate)(nil), &mrpostgres.SQLBuilderUpdate{})
}
