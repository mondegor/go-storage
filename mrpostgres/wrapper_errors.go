package mrpostgres

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mondegor/go-sysmess/errors"
)

func wrapError(err error) error {
	if err.Error() == "unexpected EOF" {
		return errors.ErrSystemStorageUnexpectedEOF.New("source", connectionName)
	}

	if e := (*pgconn.PgError)(nil); errors.As(err, &e) {
		// Severity: ERROR; Code: 42601; Message syntax error at or near "item_status"
		//           ERROR: invalid input syntax for type uuid: \"some-string\" (SQLSTATE 22P02)
		return errors.ErrInternalStorageQueryFailed.Wrap(err, "source", connectionName)
	}

	return errors.WrapInternalError(err, "failed", "source", connectionName)
}

func wrapErrorFetchDataFailed(err error) error {
	if err.Error() == "unexpected EOF" {
		return errors.ErrSystemStorageUnexpectedEOF.New("source", connectionName)
	}

	return errors.ErrInternalStorageFetchDataFailed.Wrap(err, "source", connectionName)
}

func wrapErrorCommandTag(commandTag pgconn.CommandTag, err error) error {
	if err != nil {
		return wrapError(err)
	}

	if commandTag.RowsAffected() < 1 {
		if commandTag.Update() || commandTag.Delete() {
			return errors.ErrEventStorageRowsNotAffected
		}
	}

	return nil
}
