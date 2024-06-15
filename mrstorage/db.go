package mrstorage

import "context"

type (
	// DBConnManager - менеджер соединений с БД.
	DBConnManager interface {
		Conn(ctx context.Context) DBConn
		DBTxManager
	}

	// DBTxManager - менеджер транзакций, позволяет исполнять
	// несколько независимых запросов в рамках одной транзакции.
	DBTxManager interface {
		Do(ctx context.Context, job func(ctx context.Context) error) error
	}

	// DBConn - соединение с БД с возможностью выполнения запросов.
	DBConn interface {
		Query(ctx context.Context, sql string, args ...any) (DBQueryRows, error)
		QueryRow(ctx context.Context, sql string, args ...any) DBQueryRow
		Exec(ctx context.Context, sql string, args ...any) error
	}

	// DBQueryRows - результат запроса в виде списка записей.
	DBQueryRows interface {
		Next() bool
		Scan(dest ...any) error
		Err() error
		Close()
	}

	// DBQueryRow - результат запроса состоящий из одной записи.
	DBQueryRow interface {
		Scan(dest ...any) error
	}
)
