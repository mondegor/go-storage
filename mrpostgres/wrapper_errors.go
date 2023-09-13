package mrpostgres

import (
    "github.com/jackc/pgx/v5/pgconn"
    "github.com/mondegor/go-webcore/mrcore"
)

func (c *ConnAdapter) wrapError(err error) error {
    _, ok := err.(*pgconn.PgError)

    if ok {
        // Severity: ERROR; Code: 42601; Message syntax error at or near "item_status"
        return mrcore.FactoryErrStorageQueryFailed.Caller(2).Wrap(err)
    }

    if err.Error() == "no rows in result set" {
        return mrcore.FactoryErrStorageNoRowFound.Caller(2).Wrap(err)
    }

    return mrcore.FactoryErrInternal.Caller(2).Wrap(err)
}
