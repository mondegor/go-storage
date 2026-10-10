package pgtest

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // WARNING: используется в migrate.NewWithDatabaseInstance
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrpostgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/go-storage/mrtests/internal/testenv"
)

const (
	dockerImageEnv = "MRTESTS_POSTGRES_DOCKER_IMAGE"
	dockerImage    = "postgres:18.3-alpine3.23"
	dbName         = "db_pg_test"
	dbUser         = "user_pg"
	dbPassword     = "123456_test"
	defaultSchema  = "public"
)

type (
	// Tester - вспомогательный объект для работы с тестовой БД.
	Tester struct {
		truncateCondition string
		conn              *mrpostgres.ConnAdapter
		connManager       *mrpostgres.ConnManager
	}
)

// NewTester - создаёт объект Tester.
// dbSchemas - список схем в которых будет происходить очистка таблиц,
// если не указан, то будет использоваться схема defaultSchema.
// excludedTables - список таблиц, которые будут исключены из очистки таблиц,
// имя таблицы указывается со схемой или без неё (тогда она ищется через search_path).
// Таблица миграций migratepostgres.DefaultMigrationsTable исключается всегда
// (особенности очистки исключённых таблиц см. в TruncateTables).
// Имена схем и таблиц задаются автором теста и подставляются в SQL как есть,
// поэтому передавать в них внешние данные нельзя. Имена приводятся через ::regnamespace
// и ::regclass намеренно: опечатка в имени схемы или таблицы приводит к ошибке
// при вызове TruncateTables, а не к тихому пропуску (поэтому не to_regclass).
// Переданные срезы не изменяются.
// Докер образ берётся из DockerImage.
// Соединение и контейнер освобождаются автоматически по завершении теста t,
// а методы тестера принимают t того теста (подтеста), из которого они вызываются.
func NewTester(t *testing.T, dbSchemas, excludedTables []string) *Tester {
	t.Helper()

	ctx := context.Background()
	container, err := NewContainer(
		ctx,
		DockerImage(),
		dbName,
		dbUser,
		dbPassword,
	)
	require.NoError(t, err)

	conn, err := newPostgres(ctx, container.DSN())
	if err != nil {
		_ = container.Terminate(ctx)
	}

	require.NoError(t, err)

	// ресурсы освобождаются по завершении теста t (после всех его подтестов)
	t.Cleanup(func() {
		assert.NoError(t, errors.Join(conn.Close(), container.Terminate(context.Background())))
	})

	return &Tester{
		truncateCondition: prepareTruncateCondition(dbSchemas, excludedTables),
		conn:              conn,
		connManager:       mrpostgres.NewConnManager(conn, mrlog.NopLogger()),
	}
}

// DockerImage - возвращает докер образ Postgres из переменной окружения
// MRTESTS_POSTGRES_DOCKER_IMAGE, а если она не задана, то берётся значение по умолчанию.
func DockerImage() string {
	return testenv.DockerImage(dockerImageEnv, dockerImage)
}

// ConnManager - возвращает менеджер текущего соединения с БД.
func (pt *Tester) ConnManager() *mrpostgres.ConnManager {
	return pt.connManager
}

// TruncateTables - очищает все таблицы схем dbSchemas, кроме excludedTables,
// со сбросом счётчика автоинкремента.
// Вызывается после ApplyMigrations: до них таблицы миграций ещё нет и её приведение
// к ::regclass завершается ошибкой.
// Если очищать нечего (неверно указаны схемы или исключены все таблицы),
// то тест завершается ошибкой с понятным сообщением, т.к. это ошибка настройки теста.
// WARNING: очистка выполняется с CASCADE, поэтому таблицы, ссылающиеся внешним ключом
// на очищаемые, будут очищены, даже если они указаны в excludedTables.
func (pt *Tester) TruncateTables(t *testing.T, ctx context.Context) {
	t.Helper()

	sql := fmt.Sprintf(`
		DO $do$
		DECLARE
			tables text;
		BEGIN
			SELECT string_agg(oid::regclass::text, ', ') INTO tables
			FROM pg_class
			WHERE relkind = 'r'%s;

			IF tables IS NULL THEN
				RAISE EXCEPTION 'pgtest: no tables to truncate (check dbSchemas and excludedTables)';
			END IF;

			EXECUTE 'TRUNCATE TABLE ' || tables || ' RESTART IDENTITY CASCADE';
		END $do$;`,
		pt.truncateCondition,
	)

	err := pt.conn.Exec(ctx, sql)
	require.NoError(t, err)
}

// ApplyMigrations - накатывает миграции расположенные в указанной директории.
func (pt *Tester) ApplyMigrations(t *testing.T, dirPath string) {
	t.Helper()

	pgxPool, err := pt.conn.Cli()
	require.NoError(t, err)

	db := stdlib.OpenDBFromPool(pgxPool)

	defer func() { _ = db.Close() }()

	driver, err := migratepostgres.WithInstance(db, &migratepostgres.Config{})
	require.NoError(t, err)

	dbMigrate, err := migrate.NewWithDatabaseInstance("file://"+dirPath, dbName, driver)
	require.NoError(t, err)

	defer func() { _, _ = dbMigrate.Close() }()

	err = dbMigrate.Up()
	require.NoError(t, err)
}

// ApplyFixtures - загружает данные из файлов .yml/.yaml указанной директории в БД,
// имя файла без расширения = [схема + '.'] + имя таблицы
// (без схемы таблица ищется через search_path).
// Перед добавлением данных таблица будет очищена.
func (pt *Tester) ApplyFixtures(t *testing.T, dirPath string) {
	t.Helper()

	pgxPool, err := pt.conn.Cli()
	require.NoError(t, err)

	db := stdlib.OpenDBFromPool(pgxPool)

	defer func() { _ = db.Close() }()

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(dirPath),
	)
	require.NoError(t, err)

	require.NoError(t, fixtures.Load())
}

// CountRows - возвращает количество записей указанной таблицы.
// Имя таблицы указывается со схемой или без неё (тогда она ищется через search_path)
// и подставляется в SQL как есть, поэтому передавать в него внешние данные нельзя.
func (pt *Tester) CountRows(t *testing.T, ctx context.Context, tableName string) (count int) {
	t.Helper()

	err := pt.conn.
		QueryRow(ctx, `SELECT COUNT(*) FROM `+tableName).
		Scan(&count)

	require.NoError(t, err)

	return count
}

func newPostgres(ctx context.Context, dsn string) (*mrpostgres.ConnAdapter, error) {
	conn := mrpostgres.New()
	opts := mrpostgres.Options{
		DSN: dsn,
	}

	if err := conn.Connect(ctx, opts); err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return conn, nil
}

func prepareTruncateCondition(dbSchemas, excludedTables []string) (condition string) {
	if len(dbSchemas) == 0 {
		dbSchemas = []string{defaultSchema}
	}

	condition = " AND relnamespace IN ('" + strings.Join(dbSchemas, "'::regnamespace,'") + "'::regnamespace)"

	// таблица миграций исключается всегда; срез собирается заново, чтобы не изменять переданный;
	// значения regclass сравниваются по oid, поэтому 'public.t' и 't' - одна и та же таблица
	excludedTables = append([]string{migratepostgres.DefaultMigrationsTable}, excludedTables...)
	condition += " AND oid::regclass NOT IN ('" + strings.Join(excludedTables, "'::regclass,'") + "'::regclass)"

	return condition
}
