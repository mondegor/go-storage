package mrstorage

import "context"

type (
    DbConn interface {
        Begin(ctx context.Context) (DbTransaction, error)
        DbQuery
    }

    DbTransaction interface {
        Commit(ctx context.Context) error
        Rollback(ctx context.Context) error
        DbQuery
    }

    DbQuery interface {
        Query(ctx context.Context, sql string, args ...any) (DbQueryRows, error)
        QueryRow(ctx context.Context, sql string, args ...any) DbQueryRow
        Exec(ctx context.Context, sql string, args ...any) error
    }

    DbQueryRows interface {
        Next() bool
        Scan(dest ...any) error
        Err() error
        Close()
    }

    DbQueryRow interface {
        Scan(dest ...any) error
    }
)
