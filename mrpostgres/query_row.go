package mrpostgres

import (
	"github.com/jackc/pgx/v5"
	"github.com/mondegor/go-sysmess/errors"
)

type (
	queryRow struct {
		row pgx.Row
	}
)

// Scan - comment method.
func (qr *queryRow) Scan(dest ...any) error {
	if err := qr.row.Scan(dest...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.ErrEventStorageNoRowFound
		}

		if errors.Is(err, pgx.ErrTooManyRows) {
			return errors.ErrInternalStorageFetchDataFailed.WithDetails(
				"too many rows",
				"source", connectionName,
			)
		}

		return wrapError(err)
	}

	return nil
}
