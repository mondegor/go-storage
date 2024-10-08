package mrpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the mrpostgres.SQLBuilderLimit conforms with the mrstorage.SQLBuilderLimit interface.
func TestSQLBuilderLimitImplementsSQLBuilderLimit(t *testing.T) {
	assert.Implements(t, (*mrstorage.SQLBuilderLimit)(nil), &mrpostgres.SQLBuilderLimit{})
}
