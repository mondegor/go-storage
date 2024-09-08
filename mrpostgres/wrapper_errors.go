package mrpostgres

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

func wrapError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Severity: ERROR; Code: 42601; Message syntax error at or near "item_status"
		return mrcore.ErrStorageQueryFailed.Wrap(err)
	}

	return mrcore.ErrInternal.Wrap(err)
}

func traceQuery(ctx context.Context, sql string) {
	mrlog.Ctx(ctx).
		Trace().
		Str("source", connectionName).
		Str("SQL", strings.Join(strings.Fields(sql), " ")).
		Send()
}
