package logger

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/mondegor/go-webcore/mrlog"
)

// QueryTracer - traces Query, QueryRow, and Exec.
type QueryTracer struct {
	source   string
	database string
}

// NewQueryTracer - создаёт объект QueryTracer.
func NewQueryTracer(host, port, database string) *QueryTracer {
	return &QueryTracer{
		source:   host + ":" + port,
		database: database,
	}
}

// TraceQueryStart - вызывается в начале запросов: Query, QueryRow, and Exec.
func (t *QueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	mrlog.Ctx(ctx).
		Trace().
		Str("source", t.source).
		Str("database", t.database).
		Msgf("SQL: %s ARGS: %v", strings.Join(strings.Fields(data.SQL), " "), data.Args)

	return ctx
}

// TraceQueryEnd - вызывается в конце запросов: Query, QueryRow, and Exec.
func (t *QueryTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	mrlog.Ctx(ctx).
		Trace().
		Str("source", t.source).
		Str("database", t.database).
		Msgf("CommandTag: %s; err: %v", data.CommandTag.String(), data.Err)
}
