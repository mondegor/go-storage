package mrsql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the mrsql.EntityMeta conforms with the mrstorage.MetaGetter interface.
func TestEntityMetaImplementsMetaGetter(t *testing.T) {
	assert.Implements(t, (*mrstorage.MetaGetter)(nil), &mrsql.EntityMeta{})
}
