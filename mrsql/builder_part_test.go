package mrsql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the mrsql.BuilderPart conforms with the mrstorage.SQLBuilderPart interface.
func TestBuilderPartImplementsSQLBuilderPart(t *testing.T) {
	assert.Implements(t, (*mrstorage.SQLBuilderPart)(nil), &mrsql.BuilderPart{})
}
