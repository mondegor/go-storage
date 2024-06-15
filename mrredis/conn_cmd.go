package mrredis

import (
	"context"
	"strings"
	"time"
)

// GetStruct - comment method.
func (c *ConnAdapter) GetStruct(ctx context.Context, key string, data any) error {
	c.traceCmd(ctx, "get-struct", key, data)

	if err := c.conn.Get(ctx, key).Scan(data); err != nil {
		return c.wrapError(err)
	}

	return nil
}

// SetStruct - comment method.
func (c *ConnAdapter) SetStruct(ctx context.Context, key string, data any, expiration time.Duration) error {
	c.traceCmd(ctx, "set-struct", key, data)

	if err := c.conn.Set(ctx, key, data, expiration).Err(); err != nil {
		return c.wrapError(err)
	}

	return nil
}

// Delete - comment method.
func (c *ConnAdapter) Delete(ctx context.Context, key ...string) error {
	c.traceCmd(ctx, "delete-row", strings.Join(key, ", "), nil)

	if err := c.conn.Del(ctx, key...).Err(); err != nil {
		return c.wrapError(err)
	}

	return nil
}
