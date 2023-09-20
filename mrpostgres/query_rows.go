package mrpostgres

import (
    "github.com/jackc/pgx/v5"
    "github.com/mondegor/go-webcore/mrcore"
)

type (
    queryRows struct {
        rows pgx.Rows
    }
)

func (qr *queryRows) Next() bool {
    return qr.rows.Next()
}

func (qr *queryRows) Scan(dest ...any) error {
    err := qr.rows.Scan(dest...)

    if err != nil {
        return mrcore.FactoryErrStorageFetchDataFailed.Caller(1).Wrap(err)
    }

    return nil
}

func (qr *queryRows) Err() error {
    err := qr.rows.Err()

    if err != nil {
        return mrcore.FactoryErrStorageFetchDataFailed.Caller(1).Wrap(err)
    }

    return nil
}

func (qr *queryRows) Close() {
    qr.rows.Close()
}
