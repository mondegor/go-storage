package mrpostgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-storage/mrstorage"
)

// ConnManager - менеджер транзакций.
type (
	ConnManager struct {
		conn   *ConnAdapter
		logger mrlog.Logger
	}

	ctxTransactionKey struct{}
)

// NewConnManager - создаёт объект ConnManager.
func NewConnManager(conn *ConnAdapter, logger mrlog.Logger) *ConnManager {
	return &ConnManager{
		conn:   conn,
		logger: logger,
	}
}

// Conn - возвращает соединение с PostgreSQL или транзакцию, если она была открыта.
func (m *ConnManager) Conn(ctx context.Context) mrstorage.DBConn {
	if tx, ok := ctx.Value(ctxTransactionKey{}).(*transaction); ok {
		return tx
	}

	return m.conn
}

// ConnAdapter - возвращает соединение с PostgreSQL.
func (m *ConnManager) ConnAdapter() *ConnAdapter {
	return m.conn
}

// Do - запускает задачу с запросом в транзакции.
// Пытается запустить в текущей транзакции, если ее нет, создает новую транзакцию.
func (m *ConnManager) Do(ctx context.Context, job func(ctx context.Context) error) error {
	if _, ok := ctx.Value(ctxTransactionKey{}).(*transaction); ok {
		return job(ctx)
	}

	pgxTx, err := m.conn.pool.Begin(ctx)
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

	ctx = context.WithValue(ctx, ctxTransactionKey{}, &transaction{tx: pgxTx})

	if err = job(ctx); err != nil {
		return err
	}

	if err = pgxTx.Commit(ctx); err != nil {
		return wrapError(err)
	}

	return nil
}
