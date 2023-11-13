package mrstorage

import "context"

type (
	DBConn interface {
		Begin(ctx context.Context) (DBTransaction, error)
		DBQuery
	}

	DBTransaction interface {
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
		DBQuery
	}

	DBQuery interface {
		Query(ctx context.Context, sql string, args ...any) (DBQueryRows, error)
		QueryRow(ctx context.Context, sql string, args ...any) DBQueryRow
		Exec(ctx context.Context, sql string, args ...any) error
	}

	DBQueryRows interface {
		Next() bool
		Scan(dest ...any) error
		Err() error
		Close()
	}

	DBQueryRow interface {
		Scan(dest ...any) error
	}
)
