package db

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// RowExistsChecker - формирователь запроса для проверки существования записи в таблице.
	// Поддерживает фильтрацию по полю мягкого удаления.
	RowExistsChecker[RowID any] struct {
		client          mrstorage.DBConnManager
		sqlIsExistValue string
	}
)

// NewRowExistsChecker - создаёт объект RowExistsChecker.
// Параметры:
//   - client - менеджер подключений к БД;
//   - tableName - имя таблицы для запроса;
//   - fieldKeyName - имя ключевого поля для проверки;
//   - fieldDeletedName - имя поля мягкого удаления (может быть пустым).
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

// IsExist - проверяет, существует ли запись по указанному значению поля в таблице.
// Возвращаемые значения:
//   - nil - запись найдена;
//   - ErrEventStorageNoRecordFound - запись не найдена;
//   - error - ошибка выполнения запроса;
func (re RowExistsChecker[RowID]) IsExist(ctx context.Context, id RowID) error {
	var value uint64

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
        FETCH FIRST 1 ROW ONLY;`
}
