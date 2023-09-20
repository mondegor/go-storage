package mrpostgres

import (
    "github.com/jackc/pgx/v5"
)

type (
    queryRow struct {
        row pgx.Row
        err error
    }
)

func (qr *queryRow) Scan(dest ...any) error {
    if qr.err != nil {
        return qr.err
    }

    err := qr.row.Scan(dest...)

    if err != nil {
        return wrapError(err, skipThisMethod)
    }

    return nil
}
