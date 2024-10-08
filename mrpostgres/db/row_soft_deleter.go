package db

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// RowSoftDeleter - формирователь запроса для пометки записи таблицы как удалённая.
	RowSoftDeleter[RowID any] struct {
		client           mrstorage.DBConnManager
		sqlSoftDeleteRow string
	}
)

// NewRowSoftDeleter - создаёт объект RowSoftDeleter.
func NewRowSoftDeleter[RowID any](
	client mrstorage.DBConnManager,
	tableName, fieldKeyName, fieldVersionName, fieldDeletedName string,
) RowSoftDeleter[RowID] {
	return RowSoftDeleter[RowID]{
		client:           client,
		sqlSoftDeleteRow: prepareSQLSoftDeleteRow(tableName, fieldKeyName, fieldVersionName, fieldDeletedName),
	}
}

// Delete - помечает указанную запись в качестве удалённой, если такая существует.
func (r RowSoftDeleter[RowID]) Delete(ctx context.Context, id RowID) error {
	return r.client.Conn(ctx).Exec(
		ctx,
		r.sqlSoftDeleteRow,
		id,
	)
}

func prepareSQLSoftDeleteRow(tableName, fieldKeyName, fieldVersionName, fieldDeletedName string) string {
	var set string

	if fieldVersionName != "" {
		set = fieldVersionName + ` = ` + fieldVersionName + ` + 1, `
	}

	return `
        UPDATE
            ` + tableName + `
        SET
            ` + set + fieldDeletedName + ` = NOW()
        WHERE
            ` + fieldKeyName + ` = $1 AND ` + fieldDeletedName + ` IS NULL;`
}
