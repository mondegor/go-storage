package mrpostgres

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
)

// Query - comment method.
func (c *ConnAdapter) Query(ctx context.Context, sql string, args ...any) (mrstorage.DBQueryRows, error) {
	return c.query(ctx, c.pool, sql, args...)
}

// QueryRow - comment method.
func (c *ConnAdapter) QueryRow(ctx context.Context, sql string, args ...any) mrstorage.DBQueryRow {
	return c.queryRow(ctx, c.pool, sql, args...)
}

// Exec - comment method.
func (c *ConnAdapter) Exec(ctx context.Context, sql string, args ...any) error {
	return c.exec(ctx, c.pool, sql, args...)
}
