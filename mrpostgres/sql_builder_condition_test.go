package mrpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the mrpostgres.SQLBuilderCondition conforms with the mrstorage.SQLBuilderCondition interface.
func TestSQLBuilderConditionImplementsSQLBuilderCondition(t *testing.T) {
	assert.Implements(t, (*mrstorage.SQLBuilderCondition)(nil), &mrpostgres.SQLBuilderCondition{})
}
