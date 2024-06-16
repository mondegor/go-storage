package mrpostgres

import "context"

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
func (m *TxManagerStub) Do(ctx context.Context, job func(ctx context.Context) error) error {
	return job(ctx)
}
