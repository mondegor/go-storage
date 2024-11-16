package logger

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/mondegor/go-webcore/mrlog"
)

// QueryTracer - traces Query, QueryRow, and Exec.
type QueryTracer struct {
	source string
}

// NewQueryTracer - создаёт объект QueryTracer.
func NewQueryTracer(host, port, database string) *QueryTracer {
	return &QueryTracer{
		source: host + ":" + port + "/" + database,
	}
}

// TraceQueryStart - вызывается в начале запросов: Query, QueryRow, and Exec.
func (t *QueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	const maxArgs = 16

	lenArgs := len(data.Args)
	if lenArgs > maxArgs {
		lenArgs = maxArgs
	}

	mrlog.Ctx(ctx).
		Trace().
		Str("source", t.source).
		Msgf("SQL: %s ARGS: %v", strings.Join(strings.Fields(data.SQL), " "), data.Args[:lenArgs])

	return ctx
}

// TraceQueryEnd - вызывается в конце запросов: Query, QueryRow, and Exec.
func (t *QueryTracer) TraceQueryEnd(_ context.Context, _ *pgx.Conn, _ pgx.TraceQueryEndData) {
	// mrlog.Ctx(ctx).
	//	Trace().
	//	Str("source", t.source).
	//	Msgf("CommandTag: %s; err: %v", data.CommandTag.String(), data.Err)
}
