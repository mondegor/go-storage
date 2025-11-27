package mrlock

import "github.com/mondegor/go-sysmess/mrerr"

var (
	// ErrStorageLockKeyNotObtained - lock key not obtained (ключ блокировки не удалось захватить).
	// Это вспомогательная ошибка, для неё отключено формирование стека вызовов и отправление события о её создании.
	ErrStorageLockKeyNotObtained = mrerr.NewKindInternal("lock key not obtained", mrerr.WithDisabledCaller(), mrerr.WithDisabledOnCreated())

	// ErrStorageLockKeyNotHeld - lock key not held (ключ блокировки был освобождён ранее).
	// Это вспомогательная ошибка, для неё отключено формирование стека вызовов и отправление события о её создании.
	ErrStorageLockKeyNotHeld = mrerr.NewKindInternal("lock key not held", mrerr.WithDisabledCaller(), mrerr.WithDisabledOnCreated())
)
