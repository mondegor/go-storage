package locker

import (
	"context"
	"fmt"

	"github.com/bsm/redislock"
	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mrlock"
)

// wrapError - обёртывает ошибки Redis Locker в стандартные ошибки приложения.
func (l *Adapter) wrapError(err error, key string) error {
	if errors.Is(err, redislock.ErrNotObtained) {
		return fmt.Errorf(
			"%w [source=%s, lock_key=%s]",
			mrlock.ErrLockKeyNotObtained, lockerName, key,
		)
	}

	if errors.Is(err, redislock.ErrLockNotHeld) {
		return fmt.Errorf(
			"%w [source=%s, lock_key=%s]",
			mrlock.ErrLockKeyNotHeld, lockerName, key,
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
