package part_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrpostgres/builder/part"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the builder.SQLLimitBuilder conforms with the mrstorage.SQLLimitBuilder interface.
func TestSQLLimitBuilderImplementsSQLLimitBuilder(t *testing.T) {
	assert.Implements(t, (*mrstorage.SQLLimitBuilder)(nil), &part.SQLLimitBuilder{})
}
