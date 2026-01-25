package mrlock

import (
	"github.com/mondegor/go-sysmess/errors"
)

var (
	// ErrSystemStorageLockKeyNotObtained - lock key not obtained (ключ блокировки не удалось захватить).
	ErrSystemStorageLockKeyNotObtained = errors.NewSystemProto("lock key not obtained")

	// ErrSystemStorageLockKeyNotHeld - lock key not held (ключ блокировки был освобождён ранее).
	ErrSystemStorageLockKeyNotHeld = errors.NewSystemProto("lock key not held")
)
