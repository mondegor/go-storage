package mrredis

import (
    "context"
    "strings"
    "time"
)

func (c *ConnAdapter) GetStruct(ctx context.Context, key string, data any) error {
    if err := c.conn.Get(ctx, key).Scan(data); err != nil {
        return c.wrapError(err)
    }

    c.debugCmd(ctx, "get-struct", key, data)

    return nil
}

func (c *ConnAdapter) SetStruct(ctx context.Context, key string, data any, expiration time.Duration) error {
    c.debugCmd(ctx, "set-struct", key, data)

    if err := c.conn.Set(ctx, key, data, expiration).Err(); err != nil {
        return c.wrapError(err)
    }

    return nil
}

func (c *ConnAdapter) Delete(ctx context.Context, key ...string) error {
    c.debugCmd(ctx, "delete-row", strings.Join(key, ", "), nil)

    if err := c.conn.Del(ctx, key...).Err(); err != nil {
        return c.wrapError(err)
    }

    return nil
}
