package redislocker

import (
	"context"
	"errors"

	"github.com/bsm/redislock"
	"github.com/mondegor/go-sysmess/mrerr/mr"

	"github.com/mondegor/go-storage/mrlock"
)

func (l *LockerAdapter) wrapError(err error, key string) error {
	if errors.Is(err, redislock.ErrNotObtained) {
		return mrlock.ErrStorageLockKeyNotObtained.New(
			"source", lockerName,
			"lock_key", key,
		)
	}

	if errors.Is(err, redislock.ErrLockNotHeld) {
		return mrlock.ErrStorageLockKeyNotHeld.New(
			"source", lockerName,
			"lock_key", key,
		)
	}

	return mr.ErrStorageQueryFailed.Wrap(
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
