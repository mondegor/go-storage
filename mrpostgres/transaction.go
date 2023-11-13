package mrpostgres

import (
    "context"

    "github.com/jackc/pgx/v5"
    "github.com/mondegor/go-storage/mrstorage"
)

type (
    Transaction struct {
        tx pgx.Tx
        dbExecHelper
    }
)

func (t *Transaction) Commit(ctx context.Context) error {
    return t.tx.Commit(ctx)
}

func (t *Transaction) Rollback(ctx context.Context) error {
    return t.tx.Rollback(ctx)
}

func (t *Transaction) Query(ctx context.Context, sql string, args ...any) (mrstorage.DBQueryRows, error) {
    return t.query(t.tx, skipThisMethod, ctx, sql, args...)
}

func (t *Transaction) QueryRow(ctx context.Context, sql string, args ...any) mrstorage.DBQueryRow {
    return t.queryRow(t.tx, skipThisMethod, ctx, sql, args...)
}

func (t *Transaction) Exec(ctx context.Context, sql string, args ...any) error {
    return t.exec(t.tx, skipThisMethod, ctx, sql, args...)
}
