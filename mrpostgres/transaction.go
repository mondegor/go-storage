package mrpostgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	transaction struct {
		tx pgx.Tx
	}
)

// Query - отправляет SQL запрос к БД и возвращает результат в виде списка записей.
func (t *transaction) Query(ctx context.Context, sql string, args ...any) (mrstorage.DBQueryRows, error) {
	rows, err := t.tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, wrapError(err)
	}

	return &queryRows{
		rows: rows,
	}, nil
}

// QueryRow - отправляет SQL запрос к БД и возвращает результат в виде одной записи.
func (t *transaction) QueryRow(ctx context.Context, sql string, args ...any) mrstorage.DBQueryRow {
	return &queryRow{
		row: t.tx.QueryRow(ctx, sql, args...),
	}
}

// Exec - отправляет SQL запрос к БД и исполняет его.
func (t *transaction) Exec(ctx context.Context, sql string, args ...any) error {
	return wrapErrorCommandTag(t.tx.Exec(ctx, sql, args...))
}
