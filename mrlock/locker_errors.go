package mrlock

import (
	"github.com/mondegor/go-sysmess/errors"
)

var (
	// ErrSystemStorageLockKeyNotObtained - не удалось захватить ключ блокировки.
	// Возникает, когда ключ уже заблокирован другим процессом или истёк срок предыдущей блокировки.
	ErrSystemStorageLockKeyNotObtained = errors.NewSystemProto("lock key not obtained")

	// ErrSystemStorageLockKeyNotHeld - ключ блокировки не принадлежит текущему владельцу.
	// Возникает при попытке освободить ключ, который не был заблокирован текущим процессом.
	ErrSystemStorageLockKeyNotHeld = errors.NewSystemProto("lock key not held")
)
