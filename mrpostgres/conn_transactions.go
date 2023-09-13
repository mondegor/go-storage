package mrpostgres

import (
    "context"

    "github.com/jackc/pgx/v5"
)

func (c *ConnAdapter) Begin(ctx context.Context) (pgx.Tx, error) {
    tx, err := c.pool.Begin(ctx)

    if err != nil {
        return nil, c.wrapError(err)
    }

    return tx, nil
}
