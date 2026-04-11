package mrpostgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-storage/mrstorage/txisolevel"
)

// ConnManager - менеджер подключений и транзакций PostgreSQL.
// Позволяет выполнять запросы как в рамках транзакции, так и напрямую.
// Автоматически определяет, находится ли код в контексте транзакции.
type (
	ConnManager struct {
		conn   *ConnAdapter
		logger mrlog.Logger
	}

	// ctxTxKey - ключ для хранения транзакции в контексте.
	ctxTxKey struct{}
)

// NewConnManager - создаёт объект ConnManager.
func NewConnManager(conn *ConnAdapter, logger mrlog.Logger) *ConnManager {
	return &ConnManager{
		conn:   conn,
		logger: logger,
	}
}

// Conn - возвращает соединение с PostgreSQL или активную транзакцию из контекста.
func (m *ConnManager) Conn(ctx context.Context) mrstorage.DBConn {
	if tx, ok := ctx.Value(ctxTxKey{}).(*transaction); ok {
		return tx
	}

	return m.conn
}

// ConnAdapter - возвращает адаптер подключения к PostgreSQL.
func (m *ConnManager) ConnAdapter() *ConnAdapter {
	return m.conn
}

// Do - выполняет работу в транзакции.
// Если в контексте уже есть транзакция, выполняет работу в ней.
// Иначе создаёт новую транзакцию с указанным уровнем изоляции (по умолчанию: ReadCommitted).
// Автоматически выполняет Rollback при ошибке или panic для гарантии закрытия транзакции.
// Предупреждает в логе при попытке повысить уровень изоляции во вложенной транзакции.
func (m *ConnManager) Do(ctx context.Context, job func(ctx context.Context) error, opts ...mrstorage.TxOption) error {
	o := mrstorage.TxOptions{
		IsoLevel: txisolevel.ReadCommitted, // TODO: перенести в настройки по умолчанию для менеджера соединения
	}

	for _, opt := range opts {
		opt(&o)
	}

	if tx, ok := ctx.Value(ctxTxKey{}).(*transaction); ok {
		if o.IsoLevel > tx.isoLevel {
			m.logger.Warn(
				ctx,
				"unexpected increase isolation level",
				"expected", tx.isoLevel,
				"actual", o.IsoLevel,
			)
		}

		return job(ctx)
	}

	pgxTx, err := m.conn.pool.BeginTx(ctx, mappingTxPgxOptions(o))
	if err != nil {
		return wrapError(err)
	}

	defer func() {
		// страховка от panic: Rollback всегда вызывается в конце работы функции,
		// даже в случае вызова Commit, чтобы гарантировать закрытие транзакции
		if rbErr := pgxTx.Rollback(ctx); rbErr != nil {
			if errors.Is(rbErr, pgx.ErrTxClosed) {
				return // работа в штатном режиме, транзакция зафиксирована
			}

			m.logger.Error(ctx, "call unsuccessful tx.Rollback", "error", wrapError(rbErr))

			return
		}

		m.logger.Warn(ctx, "call tx.Rollback")
	}()

	ctx = context.WithValue(
		ctx,
		ctxTxKey{},
		&transaction{
			tx:       pgxTx,
			isoLevel: o.IsoLevel,
		},
	)

	if err = job(ctx); err != nil {
		return err
	}

	if err = pgxTx.Commit(ctx); err != nil {
		return wrapError(err)
	}

	return nil
}
