package mrpostgres

import (
	"context"
)

func (c *ConnAdapter) begin(ctx context.Context) (*transaction, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return nil, wrapError(err)
	}

	return &transaction{tx: tx}, nil
}
