package mrminio_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrminio"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the mrminio.FileProvider conforms with the mrstorage.FileProvider interface.
func TestFileProviderImplementsFileProvider(t *testing.T) {
	assert.Implements(t, (*mrstorage.FileProvider)(nil), &mrminio.FileProvider{})
}
