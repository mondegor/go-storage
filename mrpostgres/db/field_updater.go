package db

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// FieldUpdater - формирователь запроса для получения и обновления значения заданного поля таблицы.
	// Включает возможность чтения текущего значения перед обновлением.
	FieldUpdater[RowID any, FieldValue any] struct {
		fetcher        FieldFetcher[RowID, FieldValue]
		sqlUpdateValue string
	}
)

// NewFieldUpdater - создаёт объект FieldUpdater.
// Параметры:
//   - client - менеджер подключений к БД;
//   - tableName - имя таблицы для запроса;
//   - fieldKeyName - имя ключевого поля для фильтрации;
//   - fieldName - имя поля для чтения/обновления;
//   - fieldDeletedName - имя поля мягкого удаления (может быть пустым).
func NewFieldUpdater[RowID, FieldValue any](
	client mrstorage.DBConnManager,
	tableName, fieldKeyName, fieldName string,
	fieldDeletedName string, // OPTIONAL: can be empty
) FieldUpdater[RowID, FieldValue] {
	return FieldUpdater[RowID, FieldValue]{
		fetcher:        NewFieldFetcher[RowID, FieldValue](client, tableName, fieldKeyName, fieldName, fieldDeletedName),
		sqlUpdateValue: prepareSQLUpdateFieldValue(tableName, fieldKeyName, fieldName, fieldDeletedName),
	}
}

// Fetch - возвращает значение поля для указанной записи в таблице.
func (re FieldUpdater[RowID, FieldValue]) Fetch(ctx context.Context, id RowID) (FieldValue, error) {
	return re.fetcher.Fetch(ctx, id)
}

// Update - обновляет значение поля указанной записи в таблице.
// Автоматически устанавливает updated_at = NOW().
func (re FieldUpdater[RowID, FieldValue]) Update(ctx context.Context, id RowID, value FieldValue) error {
	err := re.fetcher.client.Conn(ctx).Exec(
		ctx,
		re.sqlUpdateValue,
		id,
		value,
	)
	if err != nil && errors.Is(err, errors.ErrEventStorageRecordsNotAffected) {
		return errors.ErrEventStorageNoRecordFound
	}

	return err
}

func prepareSQLUpdateFieldValue(tableName, fieldKeyName, fieldName, fieldDeletedName string) string {
	var where string

	if fieldDeletedName != "" {
		where = " AND " + fieldDeletedName + " IS NULL"
	}

	return `
        UPDATE
            ` + tableName + `
        SET
			updated_at = NOW(),
            ` + fieldName + ` = $2
        WHERE
            ` + fieldKeyName + ` = $1` + where + `;`
}
