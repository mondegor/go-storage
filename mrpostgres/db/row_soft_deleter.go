package db

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// RowSoftDeleter - формирователь запроса для пометки записи таблицы как удалённой (мягкое удаление).
	// Устанавливает в поле deleted_at текущее время и увеличивает версию записи (OPTIONAL).
	RowSoftDeleter[RowID any] struct {
		client           mrstorage.DBConnManager
		sqlSoftDeleteRow string
	}
)

// NewRowSoftDeleter - создаёт объект RowSoftDeleter.
// Параметры:
//   - client - менеджер подключений к БД;
//   - tableName - имя таблицы для запроса;
//   - fieldKeyName - имя ключевого поля для фильтрации;
//   - fieldVersionName - имя поля версии для увеличения (может быть пустым);
//   - fieldDeletedName - имя поля для отметки об удалении.
func NewRowSoftDeleter[RowID any](
	client mrstorage.DBConnManager,
	tableName, fieldKeyName, fieldVersionName, fieldDeletedName string,
) RowSoftDeleter[RowID] {
	return RowSoftDeleter[RowID]{
		client:           client,
		sqlSoftDeleteRow: prepareSQLSoftDeleteRow(tableName, fieldKeyName, fieldVersionName, fieldDeletedName),
	}
}

// Delete - помечает указанную запись как удалённую, если она существует и не была удалена ранее.
func (re RowSoftDeleter[RowID]) Delete(ctx context.Context, id RowID) error {
	err := re.client.Conn(ctx).Exec(
		ctx,
		re.sqlSoftDeleteRow,
		id,
	)
	if err != nil && errors.Is(err, errors.ErrEventStorageRecordsNotAffected) {
		return errors.ErrEventStorageNoRecordFound
	}

	return err
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
