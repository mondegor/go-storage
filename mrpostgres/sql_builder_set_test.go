package mrpostgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the mrpostgres.SQLBuilderSet conforms with the mrstorage.SQLBuilderSet interface.
func TestSQLBuilderSetImplementsSQLBuilderSet(t *testing.T) {
	assert.Implements(t, (*mrstorage.SQLBuilderSet)(nil), &mrpostgres.SQLBuilderSet{})
}
