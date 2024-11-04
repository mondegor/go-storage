package part_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrpostgres/builder/part"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the builder.SQLOrderByBuilder conforms with the mrstorage.SQLOrderByBuilder interface.
func TestSQLOrderByBuilderImplementsSQLOrderByBuilder(t *testing.T) {
	assert.Implements(t, (*mrstorage.SQLOrderByBuilder)(nil), &part.SQLOrderByBuilder{})
}
