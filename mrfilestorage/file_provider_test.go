package mrfilestorage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrfilestorage"
	"github.com/mondegor/go-storage/mrstorage"
)

// Make sure the mrfilestorage.FileProvider conforms with the mrstorage.FileProvider interface.
func TestFileProviderImplementsFileProvider(t *testing.T) {
	assert.Implements(t, (*mrstorage.FileProvider)(nil), &mrfilestorage.FileProvider{})
}
