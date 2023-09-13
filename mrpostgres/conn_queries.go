package mrpostgres

import (
    "context"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
)

func (c *ConnAdapter) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
    c.debugQuery(ctx, sql)

    commandTag, err := c.pool.Exec(ctx, sql, args...)

    if err != nil {
        return commandTag, c.wrapError(err)
    }

    return commandTag, nil
}

func (c *ConnAdapter) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
    c.debugQuery(ctx, sql)

    rows, err := c.pool.Query(ctx, sql, args...)

    if err != nil {
        return nil, c.wrapError(err)
    }

    return rows, nil
}

func (c *ConnAdapter) QueryRow(ctx context.Context, sql string, args ...any) QueryRow {
    c.debugQuery(ctx, sql)

    return QueryRow{
        conn: c,
        row: c.pool.QueryRow(ctx, sql, args...),
    }
}
