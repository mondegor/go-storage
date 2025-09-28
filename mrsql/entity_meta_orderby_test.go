package mrsql_test

import (
	"testing"

	"github.com/mondegor/go-sysmess/mrtype"
	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrsql"
)

// Make sure the mrsql.EntityMetaOrderBy conforms with the mrtype.ListSorter interface.
func TestEntityMetaOrderByImplementsListSorter(t *testing.T) {
	assert.Implements(t, (*mrtype.ListSorter)(nil), &mrsql.EntityMetaOrderBy{})
}
