package db

import (
	"context"

	"golang.org/x/exp/constraints"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// TotalRowsFetcher - формирователь запроса для получения кол-ва записей в заданной таблице.
	TotalRowsFetcher[CountRows constraints.Integer] struct {
		client        mrstorage.DBConnManager
		sqlFetchTotal string
	}
)

// NewTotalRowsFetcher - создаёт объект TotalRowsFetcher.
func NewTotalRowsFetcher[CountRows constraints.Integer](client mrstorage.DBConnManager, tableName string) TotalRowsFetcher[CountRows] {
	return TotalRowsFetcher[CountRows]{
		client:        client,
		sqlFetchTotal: prepareSQLFetchTotalRows(tableName),
	}
}

// Fetch - возвращает кол-ва записей в таблице по указанному условию.
func (r TotalRowsFetcher[CountRows]) Fetch(ctx context.Context, where mrstorage.SQLPart) (CountRows, error) {
	whereStr, whereArgs := where.WithPrefix(" WHERE ").ToSQL()

	var total CountRows

	err := r.client.Conn(ctx).QueryRow(
		ctx,
		r.sqlFetchTotal+whereStr+`;`,
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
