package helper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrpostgres/builder/helper"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the helper.SQLOrderBy conforms with the mrstorage.SQLOrderByHelper interface.
func TestSQLOrderByImplementsSQLOrderByHelper(t *testing.T) {
	assert.Implements(t, (*mrstorage.SQLOrderByHelper)(nil), &helper.SQLOrderBy{})
}
