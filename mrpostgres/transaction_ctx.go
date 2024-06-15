package mrpostgres

import "context"

type (
	ctxTransactionKey struct{}
)

func withTransactionContext(ctx context.Context, tx *transaction) context.Context {
	return context.WithValue(ctx, ctxTransactionKey{}, tx)
}

func ctxTransaction(ctx context.Context) *transaction {
	if tx, ok := ctx.Value(ctxTransactionKey{}).(*transaction); ok {
		return tx
	}

	return nil
}
