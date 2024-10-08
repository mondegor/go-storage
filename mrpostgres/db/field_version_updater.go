package db

import (
	"context"

	"golang.org/x/exp/constraints"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// FieldWithVersionUpdater - формирователь запроса для получения/обновления значения заданного поля таблицы.
	// При каждом обновлении записи происходит увеличение её версии, которая сохраняется в специально заданном поле.
	FieldWithVersionUpdater[RowID any, VersionValue constraints.Integer, FieldValue any] struct {
		fetcher        FieldFetcher[RowID, FieldValue]
		sqlUpdateValue string
	}
)

// NewFieldWithVersionUpdater - создаёт объект FieldWithVersionUpdater.
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
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error.
func (r FieldWithVersionUpdater[RowID, VersionValue, FieldValue]) Fetch(ctx context.Context, id RowID) (FieldValue, error) {
	return r.fetcher.Fetch(ctx, id)
}

// Update - обновляет значение поля указанной записи в таблице и возвращает идентификатор её новой версии.
func (r FieldWithVersionUpdater[RowID, VersionValue, FieldValue]) Update(
	ctx context.Context,
	id RowID,
	version VersionValue,
	field FieldValue,
) (VersionValue, error) {
	err := r.fetcher.client.Conn(ctx).QueryRow(
		ctx,
		r.sqlUpdateValue,
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
