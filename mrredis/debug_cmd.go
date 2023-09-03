package mrredis

import (
    "context"

    "github.com/mondegor/go-core/mrcore"
)

func (c *Connection) debugCmd(ctx context.Context, command string, key string, data any) {
    mrcore.ExtractLogger(ctx).Debug("Redis: cmd=%s, key=%s, struct=%+v", command, key, data)
}
