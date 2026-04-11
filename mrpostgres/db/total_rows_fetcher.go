package db

import (
	"context"

	"golang.org/x/exp/constraints"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// TotalRowsFetcher - формирователь запроса для получения количества записей в заданной таблице.
	// Поддерживает динамические условия WHERE для фильтрации записей.
	TotalRowsFetcher[CountRows constraints.Integer] struct {
		client        mrstorage.DBConnManager
		sqlFetchTotal string
	}
)

// NewTotalRowsFetcher - создаёт объект TotalRowsFetcher.
// Параметры:
//   - client - менеджер подключений к БД;
//   - tableName - имя таблицы для подсчёта записей.
func NewTotalRowsFetcher[CountRows constraints.Integer](client mrstorage.DBConnManager, tableName string) TotalRowsFetcher[CountRows] {
	return TotalRowsFetcher[CountRows]{
		client:        client,
		sqlFetchTotal: prepareSQLFetchTotalRows(tableName),
	}
}

// Fetch - возвращает количество записей в таблице по указанному условию WHERE.
func (re TotalRowsFetcher[CountRows]) Fetch(ctx context.Context, where mrstorage.SQLPart) (CountRows, error) {
	whereStr, whereArgs := where.WithPrefix(" WHERE ").ToSQL()

	var total CountRows

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		re.sqlFetchTotal+whereStr+`;`,
		whereArgs...,
	).Scan(
		&total,
	)

	return total, err
}

func prepareSQLFetchTotalRows(tableName string) string {
	return `
        SELECT
			COUNT(*)
        FROM
            ` + tableName
}
