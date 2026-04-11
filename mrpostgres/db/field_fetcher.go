package db

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
)

// :TODO: добавить опции: fieldDeletedName, wrapperError, ...

type (
	// FieldFetcher - формирователь запроса для получения значения заданного поля таблицы.
	// Поддерживает фильтрацию по полю мягкого удаления.
	FieldFetcher[RowID, FieldValue any] struct {
		client        mrstorage.DBConnManager
		sqlFetchValue string
	}
)

// NewFieldFetcher - создаёт объект FieldFetcher.
// Параметры:
//   - client - менеджер подключений к БД;
//   - tableName - имя таблицы для запроса;
//   - fieldKeyName - имя ключевого поля для фильтрации;
//   - fieldName - имя поля для выборки значения;
//   - fieldDeletedName - имя поля мягкого удаления (может быть пустым).
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
// Возвращаемые значения:
//   - (value, nil) - запись найдена, значение получено;
//   - (zero value, ErrEventStorageNoRecordFound) - запись не найдена;
//   - (zero value, error) - ошибка выполнения запроса.
func (re FieldFetcher[RowID, FieldValue]) Fetch(ctx context.Context, id RowID) (FieldValue, error) {
	var value FieldValue

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		re.sqlFetchValue,
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
        FETCH FIRST 1 ROW ONLY;`
}
