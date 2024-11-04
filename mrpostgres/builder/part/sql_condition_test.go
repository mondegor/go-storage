package part_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrpostgres/builder/part"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the builder.SQLConditionBuilder conforms with the mrstorage.SQLConditionBuilder interface.
func TestSQLConditionBuilderImplementsSQLConditionBuilder(t *testing.T) {
	assert.Implements(t, (*mrstorage.SQLConditionBuilder)(nil), &part.SQLConditionBuilder{})
}
