package mrredislock_test

import (
	"testing"

	"github.com/mondegor/go-sysmess/mrlock"
	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrredislock"
)

// Make sure the mrredislock.LockerAdapter conforms with the mrlock.Locker interface.
func TestLockerAdapterImplementsLocker(t *testing.T) {
	assert.Implements(t, (*mrlock.Locker)(nil), &mrredislock.LockerAdapter{})
}
