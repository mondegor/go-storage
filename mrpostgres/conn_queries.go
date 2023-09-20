package mrpostgres

import (
    "context"

    "github.com/mondegor/go-storage/mrstorage"
)

func (c *ConnAdapter) Query(ctx context.Context, sql string, args ...any) (mrstorage.DbQueryRows, error) {
    return c.query(c.pool, skipThisMethod, ctx, sql, args...)
}

func (c *ConnAdapter) QueryRow(ctx context.Context, sql string, args ...any) mrstorage.DbQueryRow {
    return c.queryRow(c.pool, skipThisMethod, ctx, sql, args...)
}

func (c *ConnAdapter) Exec(ctx context.Context, sql string, args ...any) error {
    return c.exec(c.pool, skipThisMethod, ctx, sql, args...)
}

func (c *ConnAdapter) SqQuery(ctx context.Context, query mrstorage.DbSqlizer) (mrstorage.DbQueryRows, error) {
    return c.sqQuery(c.pool, skipThisMethod, ctx, query)
}

func (c *ConnAdapter) SqQueryRow(ctx context.Context, query mrstorage.DbSqlizer) mrstorage.DbQueryRow {
    return c.sqQueryRow(c.pool, skipThisMethod, ctx, query)
}

func (c *ConnAdapter) SqExec(ctx context.Context, query mrstorage.DbSqlizer) error {
    return c.sqExec(c.pool, skipThisMethod, ctx, query)
}
