package mrredis

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/redis/go-redis/v9"
)

func (c *ConnAdapter) wrapError(err error) error {
	if err == redis.Nil {
		return mrcore.FactoryErrStorageNoRowFound.Caller(1).Wrap(err)
	}

	return mrcore.FactoryErrStorageQueryFailed.Caller(1).Wrap(err)
}

func (c *ConnAdapter) debugCmd(ctx context.Context, command, key string, data any) {
	mrctx.Logger(ctx).Debug(
		"%s: cmd=%s, key=%s, struct=%#v",
		connectionName,
		command,
		key,
		data,
	)
}
