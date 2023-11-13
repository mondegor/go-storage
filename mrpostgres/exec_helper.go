package mrpostgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mondegor/go-webcore/mrcore"
)

type (
	pgxQuery interface {
		Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
		Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
		QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	}

	dbExecHelper struct {
	}
)

func (e *dbExecHelper) query(conn pgxQuery, skip int, ctx context.Context, sql string, args ...any) (*queryRows, error) {
	debugQuery(ctx, sql)

	rows, err := conn.Query(ctx, sql, args...)

	if err != nil {
		return nil, wrapError(err, skip + 1)
	}

	return &queryRows{
		rows: rows,
	}, nil
}

func (e *dbExecHelper) queryRow(conn pgxQuery, skip int, ctx context.Context, sql string, args ...any) *queryRow {
	debugQuery(ctx, sql)

	return &queryRow{
		row: conn.QueryRow(ctx, sql, args...),
	}
}

func (e *dbExecHelper) exec(conn pgxQuery, skip int, ctx context.Context, sql string, args ...any) error {
	debugQuery(ctx, sql)

	commandTag, err := conn.Exec(ctx, sql, args...)

	if err != nil {
		return wrapError(err, skip + 1)
	}

	if commandTag.Update() && commandTag.RowsAffected() < 1 {
		return mrcore.FactoryErrStorageRowsNotAffected.Caller(skip + 1).New()
	}

	return nil
}
