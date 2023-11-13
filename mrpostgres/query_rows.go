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
    if err := qr.rows.Scan(dest...); err != nil {
        return mrcore.FactoryErrStorageFetchDataFailed.Wrap(err)
    }

    return nil
}

func (qr *queryRows) Err() error {
    if err := qr.rows.Err(); err != nil {
        return mrcore.FactoryErrStorageFetchDataFailed.Wrap(err)
    }

    return nil
}

func (qr *queryRows) Close() {
    qr.rows.Close()
}
