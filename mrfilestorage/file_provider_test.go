package mrfilestorage_test

import (
	"testing"

	"github.com/mondegor/go-core/mrstorage"
	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrfilestorage"
)

// Make sure the mrfilestorage.FileProvider conforms with the mrstorage.FileProvider interface.
func TestFileProviderImplementsFileProvider(t *testing.T) {
	assert.Implements(t, (*mrstorage.FileProvider)(nil), &mrfilestorage.FileProvider{})
}
