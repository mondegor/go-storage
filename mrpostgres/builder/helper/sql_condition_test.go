package helper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrpostgres/builder/helper"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the helper.SQLCondition conforms with the mrstorage.SQLConditionHelper interface.
func TestSQLConditionImplementsSQLConditionHelper(t *testing.T) {
	assert.Implements(t, (*mrstorage.SQLConditionHelper)(nil), &helper.SQLCondition{})
}
