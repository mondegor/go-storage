package mrredis

import (
	"context"

	"github.com/mondegor/go-core/errors"
	"github.com/redis/go-redis/v9"
)

// wrapError - обёртывает ошибки Redis в стандартные ошибки приложения.
func (c *ConnAdapter) wrapError(err error) error {
	if err == redis.Nil { //nolint:errorlint
		return errors.ErrEventStorageNoRecordFound
	}

	return errors.ErrInternalStorageQueryFailed.Wrap(err, "source", connectionName)
}

// traceCmd - логирует выполняемую операцию для трассировки.
func (c *ConnAdapter) traceCmd(ctx context.Context, command, key string, data any) {
	c.tracer.Trace(
		ctx,
		"source", connectionName,
		"cmd", command,
		"key", key,
		"data", data,
	)
}
