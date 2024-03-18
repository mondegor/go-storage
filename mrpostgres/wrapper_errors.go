package mrpostgres

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

const (
	skipThisMethodFrame = 1
)

func wrapError(err error, skipFrame int) error {
	_, ok := err.(*pgconn.PgError)

	if ok {
		// Severity: ERROR; Code: 42601; Message syntax error at or near "item_status"
		return mrcore.FactoryErrStorageQueryFailed.WithCaller(skipFrame + 1).Wrap(err)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return mrcore.FactoryErrStorageNoRowFound.Wrap(err)
	}

	return mrcore.FactoryErrInternal.WithCaller(skipFrame + 1).Wrap(err)
}

func traceQuery(ctx context.Context, sql string) {
	mrlog.Ctx(ctx).
		Trace().
		Str("source", connectionName).
		Str("SQL", strings.Join(strings.Fields(sql), " ")).
		Send()
}
