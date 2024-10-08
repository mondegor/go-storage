package mrpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the mrpostgres.SQLBuilderSelect conforms with the mrstorage.SQLBuilderSelect interface.
func TestSQLBuilderSetImplementsSQLBuilderSelect(t *testing.T) {
	assert.Implements(t, (*mrstorage.SQLBuilderSelect)(nil), &mrpostgres.SQLBuilderSelect{})
}
