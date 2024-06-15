package mrredislock

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

func (l *LockerAdapter) wrapError(err error) error {
	return mrcore.ErrStorageQueryFailed.Wrap(err)
}

func (l *LockerAdapter) traceCmd(ctx context.Context, command, key string) {
	mrlog.Ctx(ctx).
		Trace().
		Str("source", lockerName).
		Str("cmd", command).
		Str("key", key).
		Send()
}
