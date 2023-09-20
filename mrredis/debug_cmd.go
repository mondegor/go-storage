package mrredis

import (
    "context"

    "github.com/mondegor/go-webcore/mrctx"
)

func (c *ConnAdapter) debugCmd(ctx context.Context, command string, key string, data any) {
    mrctx.Logger(ctx).Debug("Redis: cmd=%s, key=%s, struct=%#v", command, key, data)
}
