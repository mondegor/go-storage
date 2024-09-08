package mrpostgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/mondegor/go-webcore/mrcore"
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
			return mrcore.ErrStorageNoRowFound.Wrap(err)
		}

		if errors.Is(err, pgx.ErrTooManyRows) {
			return mrcore.ErrStorageFetchDataFailed.Wrap(err)
		}

		return wrapError(err)
	}

	return nil
}
