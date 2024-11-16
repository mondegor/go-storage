package db

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// ColumnFetcher - формирователь запроса для получения списка значений
	// заданного поля таблицы на основе фильтрации по другому заданному полю.
	ColumnFetcher[FilterValue, FieldValue any] struct {
		client         mrstorage.DBConnManager
		sqlFetchColumn string
	}
)

// NewColumnFetcher - создаёт объект ColumnFetcher.
func NewColumnFetcher[FilterValue, FieldValue any](
	client mrstorage.DBConnManager,
	tableName, fieldKeyName, columnName string,
	fieldDeletedName string, // OPTIONAL: can be empty
) ColumnFetcher[FilterValue, FieldValue] {
	return ColumnFetcher[FilterValue, FieldValue]{
		client:         client,
		sqlFetchColumn: prepareSQLFetchColumn(tableName, fieldKeyName, columnName, fieldDeletedName),
	}
}

// Fetch - возвращает список значений полей по указанному значению поля-фильтра.
func (re ColumnFetcher[FilterValue, FieldValue]) Fetch(ctx context.Context, byValue FilterValue) ([]FieldValue, error) {
	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		re.sqlFetchColumn,
		byValue,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]FieldValue, 0)

	for cursor.Next() {
		var value FieldValue

		err = cursor.Scan(
			&value,
		)
		if err != nil {
			return nil, err
		}

		rows = append(rows, value)
	}

	return rows, cursor.Err()
}

func prepareSQLFetchColumn(tableName, fieldKeyName, columnName, fieldDeletedName string) string {
	var where string

	if fieldDeletedName != "" {
		where = " AND " + fieldDeletedName + " IS NULL"
	}

	return `
        SELECT
            ` + columnName + `
        FROM
            ` + tableName + `
        WHERE
            ` + fieldKeyName + ` = $1` + where + `
        GROUP BY
            ` + columnName + `
		ORDER BY
			` + columnName + ` ASC;`
}
