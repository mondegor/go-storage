package locker

import (
	"context"

	"github.com/bsm/redislock"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlock"
)

// wrapError - обёртывает ошибки Redis Locker в стандартные ошибки приложения.
func (l *Adapter) wrapError(err error, key string) error {
	if errors.Is(err, redislock.ErrNotObtained) {
		return mrlock.ErrSystemStorageLockKeyNotObtained.New(
			"source", lockerName,
			"lock_key", key,
		)
	}

	if errors.Is(err, redislock.ErrLockNotHeld) {
		return mrlock.ErrSystemStorageLockKeyNotHeld.New(
			"source", lockerName,
			"lock_key", key,
		)
	}

	return errors.ErrInternalStorageQueryFailed.Wrap(
		err,
		"source", lockerName,
		"lock_key", key,
	)
}

// traceCmd - логирует выполняемую операцию блокировки для трассировки.
func (l *Adapter) traceCmd(ctx context.Context, command, key string) {
	l.tracer.Trace(
		ctx,
		"source", lockerName,
		"cmd", command,
		"key", key,
	)
}
