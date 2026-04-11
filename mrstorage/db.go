package mrstorage

import "context"

type (
	// DBTxManager - менеджер транзакций.
	// Позволяет выполнять несколько запросов в рамках одной транзакции.
	DBTxManager interface {
		Do(ctx context.Context, job func(ctx context.Context) error, opts ...TxOption) error
	}

	// DBConnManager - менеджер соединений с базой данных.
	// Объединяет управление соединениями и транзакциями.
	// Позволяет получать соединение как напрямую, так и в рамках транзакции.
	DBConnManager interface {
		DBTxManager

		Conn(ctx context.Context) DBConn
	}

	// DBConn - соединение с базой данных для выполнения SQL-запросов.
	// Поддерживает запросы с возвратом множества записей, одной записи и выполнение команд.
	DBConn interface {
		Query(ctx context.Context, sql string, args ...any) (DBQueryRows, error)
		QueryRow(ctx context.Context, sql string, args ...any) DBQueryRow
		Exec(ctx context.Context, sql string, args ...any) error
	}

	// DBQueryRows - результат SQL-запроса в виде списка записей.
	// Предоставляет методы для итерации по записям и их считывания.
	DBQueryRows interface {
		Next() bool
		Scan(dest ...any) error
		Err() error
		Close()
	}

	// DBQueryRow - результат SQL-запроса, состоящий из одной записи.
	DBQueryRow interface {
		Scan(dest ...any) error
	}
)
