package mrredis

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/redis/go-redis/v9"
)

func (c *ConnAdapter) wrapError(err error) error {
	if err == redis.Nil { //nolint:errorlint
		return mr.ErrStorageNoRowFound.Wrap(err)
	}

	return mr.ErrStorageQueryFailed.Wrap(err)
}

func (c *ConnAdapter) traceCmd(ctx context.Context, command, key string, data any) {
	c.tracer.Trace(
		ctx,
		"source", connectionName,
		"cmd", command,
		"key", key,
		"data", data,
	)
}
