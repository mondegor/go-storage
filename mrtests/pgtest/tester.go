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
		container         *Container
		truncateCondition string
		conn              *mrpostgres.ConnAdapter
		connManager       *mrpostgres.ConnManager
	}
)

// NewTester - создаёт объект Tester.
// dbSchemas - список схем в которых будет происходить очистка таблиц,
// если не указан, то будет использоваться схема defaultSchema.
// excludedTables - список таблиц, которые будут исключены их очистки таблиц.
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

	excludedTables = append(excludedTables, migratepostgres.DefaultMigrationsTable)

	return &Tester{
		container:         container,
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

// TruncateTables - очищает все таблицы текущей схемы со сбросом счётчика автоинкремента.
func (pt *Tester) TruncateTables(t *testing.T, ctx context.Context) {
	t.Helper()

	sql := fmt.Sprintf(`
		DO $do$
		BEGIN
			EXECUTE
				(SELECT 'TRUNCATE TABLE ' || string_agg(oid::regclass::text, ', ') || ' RESTART IDENTITY CASCADE'
				 FROM pg_class
				 WHERE relkind = 'r'%s);
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

// ApplyFixtures - загружает данные из указанной директории (имя файла = схема + '.' + имя таблицы) в БД.
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

// CountRows - возвращает количество записей указанной таблицы находящейся в текущей схеме.
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
		dbSchemas = append(dbSchemas, defaultSchema)
	}

	condition = " AND relnamespace IN ('" + strings.Join(dbSchemas, "'::regnamespace,'") + "'::regnamespace)"

	if len(excludedTables) > 0 {
		prefix := defaultSchema + "."

		// публичная схема срезается у всех таблиц, иначе условие работать не будет правильно
		for i := range excludedTables {
			excludedTables[i] = strings.TrimPrefix(excludedTables[i], prefix)
		}

		condition += " AND oid::regclass NOT IN ('" + strings.Join(excludedTables, "'::regclass,'") + "'::regclass)"
	}

	return condition
}
