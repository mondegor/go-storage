package mrpostgres

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
)

func (c *ConnAdapter) Begin(ctx context.Context) (mrstorage.DBTransaction, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return nil, wrapError(err, skipThisMethodFrame)
	}

	return &Transaction{tx: tx}, nil
}
