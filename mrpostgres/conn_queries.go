package mrpostgres

import (
    "context"

    "github.com/mondegor/go-storage/mrstorage"
)

func (c *ConnAdapter) Query(ctx context.Context, sql string, args ...any) (mrstorage.DBQueryRows, error) {
    return c.query(c.pool, skipThisMethod, ctx, sql, args...)
}

func (c *ConnAdapter) QueryRow(ctx context.Context, sql string, args ...any) mrstorage.DBQueryRow {
    return c.queryRow(c.pool, skipThisMethod, ctx, sql, args...)
}

func (c *ConnAdapter) Exec(ctx context.Context, sql string, args ...any) error {
    return c.exec(c.pool, skipThisMethod, ctx, sql, args...)
}
