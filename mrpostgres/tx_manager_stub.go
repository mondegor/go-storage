package mrpostgres

import "context"

type (
	// TxManagerStub - Фиктивный менеджер транзакций, который
	// запускает только переданную работу без открытия транзакции.
	TxManagerStub struct{}
)

// NewTxManagerStub - Создаёт фиктивный менеджер транзакций.
func NewTxManagerStub() *TxManagerStub {
	return &TxManagerStub{}
}

// Do - запускает работу без использования транзакции.
func (m *TxManagerStub) Do(ctx context.Context, job func(ctx context.Context) error) error {
	return job(ctx)
}
