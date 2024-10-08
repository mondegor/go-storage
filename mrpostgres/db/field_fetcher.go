package db

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// FieldFetcher - формирователь запроса для получения значения заданного поля таблицы.
	FieldFetcher[RowID, FieldValue any] struct {
		client        mrstorage.DBConnManager
		sqlFetchValue string
	}
)

// NewFieldFetcher - создаёт объект FieldFetcher.
func NewFieldFetcher[RowID, FieldValue any](
	client mrstorage.DBConnManager,
	tableName, fieldKeyName, fieldName string,
	fieldDeletedName string, // OPTIONAL: can be empty
) FieldFetcher[RowID, FieldValue] {
	return FieldFetcher[RowID, FieldValue]{
		client:        client,
		sqlFetchValue: prepareSQLFetchFieldValue(tableName, fieldKeyName, fieldName, fieldDeletedName),
	}
}

// Fetch - возвращает значение поля для указанной записи в таблице.
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error.
func (r FieldFetcher[RowID, FieldValue]) Fetch(ctx context.Context, id RowID) (FieldValue, error) {
	var value FieldValue

	err := r.client.Conn(ctx).QueryRow(
		ctx,
		r.sqlFetchValue,
		id,
	).Scan(
		&value,
	)

	return value, err
}

func prepareSQLFetchFieldValue(tableName, fieldKeyName, fieldName, fieldDeletedName string) string {
	var where string

	if fieldDeletedName != "" {
		where = " AND " + fieldDeletedName + " IS NULL"
	}

	return `
        SELECT
            ` + fieldName + `
        FROM
            ` + tableName + `
        WHERE
            ` + fieldKeyName + ` = $1` + where + `
        LIMIT 1;`
}
