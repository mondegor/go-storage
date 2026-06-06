package locker_test

import (
	"testing"

	"github.com/mondegor/go-sysmess/mrlock"
	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrredis/locker"
)

// Make sure the locker.Adapter conforms with the mrlock.Locker interface.
func TestAdapterImplementsLocker(t *testing.T) {
	assert.Implements(t, (*mrlock.Locker)(nil), &locker.Adapter{})
}
