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

	dbExecHelper struct{}
)

func (e *dbExecHelper) query(ctx context.Context, conn pgxQuery, sql string, args ...any) (*queryRows, error) {
	traceQuery(ctx, sql)

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, wrapError(err)
	}

	return &queryRows{
		rows: rows,
	}, nil
}

func (e *dbExecHelper) queryRow(ctx context.Context, conn pgxQuery, sql string, args ...any) *queryRow {
	traceQuery(ctx, sql)

	return &queryRow{
		row: conn.QueryRow(ctx, sql, args...),
	}
}

func (e *dbExecHelper) exec(ctx context.Context, conn pgxQuery, sql string, args ...any) error {
	traceQuery(ctx, sql)

	commandTag, err := conn.Exec(ctx, sql, args...)
	if err != nil {
		return wrapError(err)
	}

	if commandTag.RowsAffected() < 1 {
		if commandTag.Insert() || commandTag.Update() || commandTag.Delete() {
			return mrcore.ErrStorageRowsNotAffected.New()
		}
	}

	return nil
}
