package mutexlocker_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrlock"
	"github.com/mondegor/go-storage/mrlock/mutexlocker"
)

// Make sure the Locker conforms with the mrlock.Locker interface.
func TestLockerImplementsLocker(t *testing.T) {
	assert.Implements(t, (*mrlock.Locker)(nil), &mutexlocker.Locker{})
}
