package mrredislock

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

func (l *lockerAdapter) wrapError(err error) error {
	return mrcore.FactoryErrStorageQueryFailed.WithCaller(1).Wrap(err)
}

func (l *lockerAdapter) traceCmd(ctx context.Context, command, key string) {
	mrlog.Ctx(ctx).
		Trace().
		Str("source", lockerName).
		Str("cmd", command).
		Str("key", key).
		Send()
}
