package mrpostgres

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mondegor/go-sysmess/errors"
)

func wrapError(err error) error {
	if ok, wrappedErr := wrapPgError(err); ok {
		if wrappedErr != nil {
			return wrappedErr
		}

		return errors.ErrInternalStorageQueryFailed.Wrap(err, "source", connectionName)
	}

	return errors.WrapInternalError(err, "failed", "source", connectionName)
}

func wrapPgError(err error) (ok bool, wrappedErr error) {
	if e := (*pgconn.PgError)(nil); errors.As(err, &e) {
		// Code: 23505 duplicate key value violates unique constraint
		if e.Code == "23505" {
			return true, errors.ErrInternalStorageDuplicateKeyViolation.Wrap(err)
		}

		return true, nil
	}

	if err.Error() == "unexpected EOF" {
		return true, errors.ErrSystemStorageUnexpectedEOF.New("source", connectionName)
	}

	return false, nil
}

func wrapErrorFetchDataFailed(err error) error {
	if _, wrappedErr := wrapPgError(err); wrappedErr != nil {
		return wrappedErr
	}

	return errors.ErrInternalStorageFetchDataFailed.Wrap(err, "source", connectionName)
}

func wrapErrorCommandTag(commandTag pgconn.CommandTag, err error) error {
	if err != nil {
		return wrapError(err)
	}

	if commandTag.RowsAffected() < 1 {
		if commandTag.Update() || commandTag.Delete() {
			return errors.ErrEventStorageRecordsNotAffected
		}
	}

	return nil
}
