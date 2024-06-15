package mrredis

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/redis/go-redis/v9"
)

func (c *ConnAdapter) wrapError(err error) error {
	if err == redis.Nil { //nolint:errorlint
		return mrcore.ErrStorageNoRowFound.Wrap(err)
	}

	return mrcore.ErrStorageQueryFailed.Wrap(err)
}

func (c *ConnAdapter) traceCmd(ctx context.Context, command, key string, data any) {
	mrlog.Ctx(ctx).
		Trace().
		Str("source", connectionName).
		Str("cmd", command).
		Str("key", key).
		Any("data", data).
		Send()
}
