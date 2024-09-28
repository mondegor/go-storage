package mrpostgres

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/go-storage/mrstorage"
)

// ConnManager - менеджер транзакций.
type (
	ConnManager struct {
		conn *ConnAdapter
	}
)

// NewConnManager - создаёт объект ConnManager.
func NewConnManager(conn *ConnAdapter) *ConnManager {
	return &ConnManager{
		conn: conn,
	}
}

// Conn - возвращает соединение с PostgreSQL или транзакцию, если она была открыта.
func (m *ConnManager) Conn(ctx context.Context) mrstorage.DBConn {
	if tx := ctxTransaction(ctx); tx != nil {
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
	if tx := ctxTransaction(ctx); tx != nil {
		return job(ctx)
	}

	return m.do(ctx, job)
}

func (m *ConnManager) do(ctx context.Context, job func(ctx context.Context) error) error {
	tx, err := m.conn.begin(ctx)
	if err != nil {
		return err
	}

	ctx = withTransactionContext(ctx, tx)

	if err = job(ctx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			mrlog.Ctx(ctx).Error().Err(err).Msg("before error in tx.Rollback")
			err = rbErr
		}

		return err
	}

	return tx.Commit(ctx)
}
