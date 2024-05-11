package mrredis

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/redis/go-redis/v9"
)

func (c *ConnAdapter) wrapError(err error) error {
	const skipFrame = 1

	if err == redis.Nil {
		return mrcore.FactoryErrStorageNoRowFound.WithSkipFrame(skipFrame).Wrap(err)
	}

	return mrcore.FactoryErrStorageQueryFailed.WithSkipFrame(skipFrame).Wrap(err)
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
