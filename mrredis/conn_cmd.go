package mrredis

import (
	"context"
	"strings"
	"time"
)

// GetStruct - получает значение из Redis по ключу и десериализует его в переданную структуру.
// Параметр data должен быть указателем на структуру, в которую будет записан результат.
func (c *ConnAdapter) GetStruct(ctx context.Context, key string, data any) error {
	c.traceCmd(ctx, "get-struct", key, data)

	if err := c.conn.Get(ctx, key).Scan(data); err != nil {
		return c.wrapError(err)
	}

	return nil
}

// SetStruct - сохраняет структуру в Redis с указанным ключом и временем жизни.
// Параметр data сериализуется перед сохранением (используется встроенная сериализация go-redis).
// Если expiration равен 0, ключ будет храниться бессрочно.
func (c *ConnAdapter) SetStruct(ctx context.Context, key string, data any, expiration time.Duration) error {
	c.traceCmd(ctx, "set-struct", key, data)

	if err := c.conn.Set(ctx, key, data, expiration).Err(); err != nil {
		return c.wrapError(err)
	}

	return nil
}

// Delete - удаляет один или несколько ключей из Redis.
// Принимает переменное количество ключей для удаления.
func (c *ConnAdapter) Delete(ctx context.Context, key ...string) error {
	c.traceCmd(ctx, "delete-row", strings.Join(key, ", "), nil)

	if err := c.conn.Del(ctx, key...).Err(); err != nil {
		return c.wrapError(err)
	}

	return nil
}
