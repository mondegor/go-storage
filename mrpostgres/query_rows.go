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

// Next - comment method.
func (qr *queryRows) Next() bool {
	return qr.rows.Next()
}

// Scan - comment method.
func (qr *queryRows) Scan(dest ...any) error {
	if err := qr.rows.Scan(dest...); err != nil {
		return mrcore.ErrStorageFetchDataFailed.Wrap(err)
	}

	return nil
}

// Err - comment method.
func (qr *queryRows) Err() error {
	if err := qr.rows.Err(); err != nil {
		return mrcore.ErrStorageFetchDataFailed.Wrap(err)
	}

	return nil
}

// Close - comment method.
func (qr *queryRows) Close() {
	qr.rows.Close()
}
