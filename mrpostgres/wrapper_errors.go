package mrpostgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mondegor/go-webcore/mrcore"
)

func wrapError(err error) error {
	if err.Error() == "unexpected EOF" {
		return mrcore.ErrInternalUnexpectedEOF.New()
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Severity: ERROR; Code: 42601; Message syntax error at or near "item_status"
		//           ERROR: invalid input syntax for type uuid: \"some-string\" (SQLSTATE 22P02)
		return mrcore.ErrStorageQueryFailed.Wrap(err)
	}

	return mrcore.ErrInternal.Wrap(err)
}

func wrapErrorFetchDataFailed(err error) error {
	if err.Error() == "unexpected EOF" {
		return mrcore.ErrInternalUnexpectedEOF.New()
	}

	return mrcore.ErrStorageFetchDataFailed.Wrap(err)
}
