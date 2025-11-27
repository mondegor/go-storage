package redislocker_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrlock"
	"github.com/mondegor/go-storage/mrlock/redislocker"
)

// Make sure the mrredislock.LockerAdapter conforms with the mrlock.Locker interface.
func TestLockerAdapterImplementsLocker(t *testing.T) {
	assert.Implements(t, (*mrlock.Locker)(nil), &redislocker.LockerAdapter{})
}
