package db

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// RowExistsChecker - формирователь запроса для проверки существования записи для заданного поля таблицы.
	RowExistsChecker[RowID any] struct {
		client          mrstorage.DBConnManager
		sqlIsExistValue string
	}
)

// NewRowExistsChecker - создаёт объект RowExistsChecker.
func NewRowExistsChecker[RowID any](
	client mrstorage.DBConnManager,
	tableName, fieldKeyName string,
	fieldDeletedName string, // OPTIONAL: can be empty
) RowExistsChecker[RowID] {
	return RowExistsChecker[RowID]{
		client:          client,
		sqlIsExistValue: prepareSQLCheckRowExists(tableName, fieldKeyName, fieldDeletedName),
	}
}

// IsExist - проверяет существование записи по указанному значению поля в таблице.
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re RowExistsChecker[RowID]) IsExist(ctx context.Context, id RowID) error {
	var value uint32

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		re.sqlIsExistValue,
		id,
	).Scan(
		&value,
	)

	return err
}

func prepareSQLCheckRowExists(tableName, fieldKeyName, fieldDeletedName string) string {
	var where string

	if fieldDeletedName != "" {
		where = " AND " + fieldDeletedName + " IS NULL"
	}

	return `
        SELECT
            1
        FROM
            ` + tableName + `
        WHERE
            ` + fieldKeyName + ` = $1` + where + `
        LIMIT 1;`
}
