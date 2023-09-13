package mrpostgres

import (
    "context"

    "github.com/jackc/pgx/v5"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrcore"
)

func (c *ConnAdapter) SqQuery(ctx context.Context, query mrstorage.Sqlizer) (pgx.Rows, error) {
    sql, args, err := query.ToSql()

    if err != nil {
        return nil, mrcore.FactoryErrInternal.Caller(1).Wrap(err)
    }

    c.debugQuery(ctx, sql)

    rows, err := c.pool.Query(ctx, sql, args...)

    if err != nil {
        return nil, c.wrapError(err)
    }

    return rows, nil
}

func (c *ConnAdapter) SqQueryRow(ctx context.Context, query mrstorage.Sqlizer) QueryRow {
    sql, args, err := query.ToSql()

    if err != nil {
        return QueryRow{err: err}
    }

    c.debugQuery(ctx, sql)

    return QueryRow{
        conn: c,
        row: c.pool.QueryRow(ctx, sql, args...),
    }
}

func (c *ConnAdapter) SqUpdate(ctx context.Context, query mrstorage.Sqlizer) error {
    sql, args, err := query.ToSql()

    if err != nil {
        return mrcore.FactoryErrInternal.Caller(1).Wrap(err)
    }

    c.debugQuery(ctx, sql)

    commandTag, err := c.pool.Exec(ctx, sql, args...)

    if err != nil {
        return c.wrapError(err)
    }

    if commandTag.RowsAffected() < 1 {
        return mrcore.FactoryErrStorageRowsNotAffected.Caller(1).New()
    }

    return nil
}
