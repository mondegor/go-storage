package mrpostgres

import (
    "context"

    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrcore"
)

func (e *dbExecHelper) sqQuery(conn pgxQuery, skip int, ctx context.Context, query mrstorage.DbSqlizer) (*queryRows, error) {
    sql, args, err := e.parseSql(skip + 1, ctx, query)

    if err != nil {
        return nil, err
    }

    return e.query(conn, skip + 1, ctx, sql, args...)
}

func (e *dbExecHelper) sqQueryRow(conn pgxQuery, skip int, ctx context.Context, query mrstorage.DbSqlizer) *queryRow {
    sql, args, err := e.parseSql(skip + 1, ctx, query)

    return &queryRow{
        row: conn.QueryRow(ctx, sql, args...),
        err: err,
    }
}

func (e *dbExecHelper) sqExec(conn pgxQuery, skip int, ctx context.Context, query mrstorage.DbSqlizer) error {
    sql, args, err := e.parseSql(skip + 1, ctx, query)

    if err != nil {
        return err
    }

    return e.exec(conn, skip + 1, ctx, sql, args...)
}

func (e *dbExecHelper) parseSql(skip int, ctx context.Context, query mrstorage.DbSqlizer) (string, []any, error) {
    sql, args, err := query.ToSql()

    if err != nil {
        return "", nil, mrcore.FactoryErrInternalParseData.Caller(skip + 1).Wrap(err)
    }

    debugQuery(ctx, sql)

    return sql, args, err
}
