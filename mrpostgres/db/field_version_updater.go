package db

import (
	"context"

	"golang.org/x/exp/constraints"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// FieldWithVersionUpdater - формирователь запроса для получения и обновления значения поля с контролем версий.
	// При каждом обновлении увеличивает версию записи для предотвращения конфликтов параллельного изменения.
	// Использует оптимистичную блокировку (optimistic locking).
	FieldWithVersionUpdater[RowID any, VersionValue constraints.Integer, FieldValue any] struct {
		fetcher        FieldFetcher[RowID, FieldValue]
		sqlUpdateValue string
	}
)

// NewFieldWithVersionUpdater - создаёт объект FieldWithVersionUpdater.
// Параметры:
//   - client - менеджер подключений к БД;
//   - tableName - имя таблицы для запроса;
//   - fieldKeyName - имя ключевого поля для фильтрации;
//   - fieldVersionName - имя поля для хранения версии записи;
//   - fieldName - имя поля для чтения/обновления;
//   - fieldDeletedName - имя поля мягкого удаления (может быть пустым).
func NewFieldWithVersionUpdater[RowID any, VersionValue constraints.Integer, FieldValue any](
	client mrstorage.DBConnManager,
	tableName, fieldKeyName, fieldVersionName, fieldName string,
	fieldDeletedName string, // OPTIONAL: can be empty
) FieldWithVersionUpdater[RowID, VersionValue, FieldValue] {
	return FieldWithVersionUpdater[RowID, VersionValue, FieldValue]{
		fetcher:        NewFieldFetcher[RowID, FieldValue](client, tableName, fieldKeyName, fieldName, fieldDeletedName),
		sqlUpdateValue: prepareSQLUpdateFieldValueWithVersion(tableName, fieldKeyName, fieldVersionName, fieldName, fieldDeletedName),
	}
}

// Fetch - возвращает значение поля для указанной записи в таблице.
// Делегирует вызов внутреннему FieldFetcher.
func (re FieldWithVersionUpdater[RowID, VersionValue, FieldValue]) Fetch(ctx context.Context, id RowID) (FieldValue, error) {
	return re.fetcher.Fetch(ctx, id)
}

// Update - обновляет значение поля указанной записи в таблице с проверкой версии.
// Увеличивает версию записи на 1 и обновляет updated_at = NOW().
// Если текущая версия не совпадает с ожидаемой, запись не будет обновлена (optimistic locking).
// Возвращает новую версию записи или ошибку, если запись не найдена или версия устарела.
func (re FieldWithVersionUpdater[RowID, VersionValue, FieldValue]) Update(
	ctx context.Context,
	id RowID,
	version VersionValue,
	field FieldValue,
) (VersionValue, error) {
	err := re.fetcher.client.Conn(ctx).QueryRow(
		ctx,
		re.sqlUpdateValue,
		id,
		version,
		field,
	).Scan(
		&version,
	)

	return version, err
}

func prepareSQLUpdateFieldValueWithVersion(tableName, fieldKeyName, fieldVersionName, fieldName, fieldDeletedName string) string {
	var where string

	if fieldDeletedName != "" {
		where = " AND " + fieldDeletedName + " IS NULL"
	}

	return `
        UPDATE
            ` + tableName + `
        SET
            ` + fieldVersionName + ` = ` + fieldVersionName + ` + 1,
			updated_at = NOW(),
            ` + fieldName + ` = $3
        WHERE
            ` + fieldKeyName + ` = $1 AND ` + fieldVersionName + ` = $2` + where + `
		RETURNING
			` + fieldVersionName + `;`
}
