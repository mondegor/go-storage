package mrpostgres

import (
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-sysmess/mrerr"

    "github.com/jackc/pgx/v5/pgconn"
)

func (c *Connection) wrapError(err error) error {
    _, ok := err.(*pgconn.PgError)

    if ok {
        // Severity: ERROR; Code: 42601; Message syntax error at or near "item_status"
        return mrstorage.ErrFactoryQueryFailed.Caller(2).Wrap(err)
    }

    if err.Error() == "no rows in result set" {
        return mrstorage.ErrFactoryNoRowFound.Caller(2).Wrap(err)
    }

    return mrerr.ErrFactoryInternal.Caller(2).Wrap(err)
}
