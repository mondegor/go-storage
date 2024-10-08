package mrsql_test

import (
	"testing"

	"github.com/mondegor/go-webcore/mrview"
	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrsql"
)

// Make sure the mrsql.EntityMetaOrderBy conforms with the mrview.ListSorter interface.
func TestEntityMetaOrderByImplementsListSorter(t *testing.T) {
	assert.Implements(t, (*mrview.ListSorter)(nil), &mrsql.EntityMetaOrderBy{})
}
