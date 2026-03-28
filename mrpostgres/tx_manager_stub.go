package mrpostgres

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// TxManagerStub - фиктивный менеджер транзакций, который
	// запускает только переданную работу без открытия транзакции.
	TxManagerStub struct{}
)

// NewTxManagerStub - создаёт объект TxManagerStub.
func NewTxManagerStub() *TxManagerStub {
	return &TxManagerStub{}
}

// Do - запускает работу без использования транзакции.
func (m *TxManagerStub) Do(ctx context.Context, job func(ctx context.Context) error, _ ...mrstorage.TxOption) error {
	return job(ctx)
}
