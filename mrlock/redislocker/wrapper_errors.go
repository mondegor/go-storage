package redislocker

import (
	"context"

	"github.com/bsm/redislock"
	"github.com/mondegor/go-sysmess/errors"

	"github.com/mondegor/go-storage/mrlock"
)

func (l *LockerAdapter) wrapError(err error, key string) error {
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

func (l *LockerAdapter) traceCmd(ctx context.Context, command, key string) {
	l.tracer.Trace(
		ctx,
		"source", lockerName,
		"cmd", command,
		"key", key,
	)
}
