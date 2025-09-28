package mrpostgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mondegor/go-sysmess/mrerr/mr"
)

func wrapError(err error) error {
	if err.Error() == "unexpected EOF" {
		return mr.ErrInternalUnexpectedEOF.New()
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Severity: ERROR; Code: 42601; Message syntax error at or near "item_status"
		//           ERROR: invalid input syntax for type uuid: \"some-string\" (SQLSTATE 22P02)
		return mr.ErrStorageQueryFailed.Wrap(err)
	}

	return mr.ErrInternal.Wrap(err)
}

func wrapErrorFetchDataFailed(err error) error {
	if err.Error() == "unexpected EOF" {
		return mr.ErrInternalUnexpectedEOF.New()
	}

	return mr.ErrStorageFetchDataFailed.Wrap(err)
}

func wrapErrorCommandTag(commandTag pgconn.CommandTag, err error) error {
	if err != nil {
		return wrapError(err)
	}

	if commandTag.RowsAffected() < 1 {
		if commandTag.Insert() || commandTag.Update() || commandTag.Delete() {
			return mr.ErrStorageRowsNotAffected.New()
		}
	}

	return nil
}
