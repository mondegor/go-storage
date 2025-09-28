package mrredislock

import (
	"context"
	"errors"

	"github.com/bsm/redislock"
	"github.com/mondegor/go-sysmess/mrerr/mr"
)

func (l *LockerAdapter) wrapError(err error, key string) error {
	args := []any{"source", lockerName, "key", key}

	if errors.Is(err, redislock.ErrNotObtained) {
		return mr.ErrStorageLockKeyNotCaptured.Wrap(err, args...)
	}

	if errors.Is(err, redislock.ErrLockNotHeld) {
		err = mr.ErrStorageLockKeyNotHeld.Wrap(err, args...)
	}

	return mr.ErrStorageQueryFailed.Wrap(err, args...)
}

func (l *LockerAdapter) traceCmd(ctx context.Context, command, key string) {
	l.tracer.Trace(
		ctx,
		"source", lockerName,
		"cmd", command,
		"key", key,
	)
}
