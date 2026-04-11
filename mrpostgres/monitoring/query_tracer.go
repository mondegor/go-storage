package monitoring

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/mondegor/go-sysmess/mrtrace"
)

// QueryTracer - трассировщик SQL-запросов для библиотеки pgx.
// Отслеживает начало выполнения запросов Query, QueryRow и Exec.
type QueryTracer struct {
	tracer mrtrace.Tracer
	source string
}

// NewQueryTracer - создаёт объект QueryTracer для трассировки SQL-запросов.
// Параметры:
//   - host - адрес сервера PostgreSQL;
//   - port - порт сервера PostgreSQL;
//   - database - имя БД;
//   - tracer - трассировщик для логирования запросов.
func NewQueryTracer(host, port, database string, tracer mrtrace.Tracer) *QueryTracer {
	return &QueryTracer{
		tracer: tracer,
		source: host + ":" + port + "/" + database,
	}
}

// TraceQueryStart - вызывается библиотекой pgx в начале выполнения SQL-запроса.
// Логирует SQL-запрос (с нормализованными пробелами) и первые 16 аргументов.
func (t *QueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	const maxArgs = 16

	lenArgs := len(data.Args)
	if lenArgs > maxArgs {
		lenArgs = maxArgs
	}

	t.tracer.Trace(
		ctx,
		"source", t.source,
		"sql", strings.Join(strings.Fields(data.SQL), " "),
		"args", data.Args[:lenArgs],
	)

	return ctx
}

// TraceQueryEnd - вызывается библиотекой pgx в конце выполнения SQL-запроса.
func (t *QueryTracer) TraceQueryEnd(_ context.Context, _ *pgx.Conn, _ pgx.TraceQueryEndData) {
	// mrlog.Ctx(ctx).
	//	Trace().
	//	Str("source", t.source).
	//	Msgf("CommandTag: %s; err: %v", data.CommandTag.String(), data.Err)
}
