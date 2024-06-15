package mrpostgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	transaction struct {
		tx pgx.Tx
		dbExecHelper
	}
)

// Commit - comment method.
func (t *transaction) Commit(ctx context.Context) error {
	if err := t.tx.Commit(ctx); err != nil {
		return wrapError(err)
	}

	return nil
}

// Rollback - comment method.
func (t *transaction) Rollback(ctx context.Context) error {
	if err := t.tx.Rollback(ctx); err != nil {
		return wrapError(err)
	}

	return nil
}

// Query - comment method.
func (t *transaction) Query(ctx context.Context, sql string, args ...any) (mrstorage.DBQueryRows, error) {
	return t.query(ctx, t.tx, sql, args...)
}

// QueryRow - comment method.
func (t *transaction) QueryRow(ctx context.Context, sql string, args ...any) mrstorage.DBQueryRow {
	return t.queryRow(ctx, t.tx, sql, args...)
}

// Exec - comment method.
func (t *transaction) Exec(ctx context.Context, sql string, args ...any) error {
	return t.exec(ctx, t.tx, sql, args...)
}
