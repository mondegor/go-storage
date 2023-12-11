package mrredislock

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
)

func (l *lockerAdapter) wrapError(err error) error {
	return mrcore.FactoryErrStorageQueryFailed.Caller(1).Wrap(err)
}

func (l *lockerAdapter) debugCmd(ctx context.Context, command, key string) {
	mrctx.Logger(ctx).Debug(
		"%s: cmd=%s, key=%s",
		lockerName,
		command,
		key,
	)
}
