package mrpostgres

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
)

const (
	skipThisMethod = 1
)

func wrapError(err error, skip int) error {
	_, ok := err.(*pgconn.PgError)

	if ok {
		// Severity: ERROR; Code: 42601; Message syntax error at or near "item_status"
		return mrcore.FactoryErrStorageQueryFailed.Caller(skip + 1).Wrap(err)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return mrcore.FactoryErrStorageNoRowFound.Caller(skip + 1).Wrap(err)
	}

	return mrcore.FactoryErrInternal.Caller(skip + 1).Wrap(err)
}

func debugQuery(ctx context.Context, sql string) {
	mrctx.Logger(ctx).Debug(
		connectionName + " SQL: " + strings.Join(strings.Fields(sql), " "),
	)
}
