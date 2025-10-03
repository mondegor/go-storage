package infra

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // WARNING: используется в migrate.NewWithDatabaseInstance
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mondegor/go-sysmess/mrlog/slog/nopslog"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrtests/helpers"
)

const (
	postgresDockerImage   = "p/postgres:15.7"
	postgresDB            = "db_pg_test"
	postgresUser          = "user_pg"
	postgresPassword      = "123456_test"
	postgresDefaultSchema = "public"
)

type (
	// PostgresTester - вспомогательный объект для работы с тестовой БД.
	PostgresTester struct {
		ownerT            *testing.T
		container         *helpers.PostgresContainer
		truncateCondition string
		conn              *mrpostgres.ConnAdapter
		connManager       *mrpostgres.ConnManager
	}
)

// NewPostgresTester - создаёт объект PostgresTester.
// dbSchemas - список схем в которых будет происходить очистка таблиц,
// если не указан, то будет использоваться схема postgresDefaultSchema.
// excludedTables - список таблиц, которые будут исключены их очистки таблиц.
func NewPostgresTester(t *testing.T, dbSchemas, excludedTables []string) *PostgresTester {
	t.Helper()

	ctx := context.Background()
	container, err := helpers.NewPostgresContainer(
		ctx,
		postgresDockerImage,
		postgresDB,
		postgresUser,
		postgresPassword,
	)
	require.NoError(t, err)

	conn, err := newPostgres(ctx, container.DSN())
	require.NoError(t, err)

	excludedTables = append(excludedTables, postgres.DefaultMigrationsTable)

	return &PostgresTester{
		ownerT:            t,
		container:         container,
		truncateCondition: prepareTruncateCondition(dbSchemas, excludedTables),
		conn:              conn,
		connManager:       mrpostgres.NewConnManager(conn, nopslog.New()),
	}
}

// ConnManager - возвращает менеджер текущего соединения с БД.
func (t *PostgresTester) ConnManager() *mrpostgres.ConnManager {
	t.ownerT.Helper()

	return t.connManager
}

// TruncateTables - очищает все таблицы текущей схемы со сбросом счётчика автоинкремента.
func (t *PostgresTester) TruncateTables(ctx context.Context) {
	t.ownerT.Helper()

	sql := fmt.Sprintf(`
		DO $do$
		BEGIN
			EXECUTE
				(SELECT 'TRUNCATE TABLE ' || string_agg(oid::regclass::text, ', ') || ' RESTART IDENTITY CASCADE'
				 FROM pg_class
				 WHERE relkind = 'r'%s);
		END $do$;`,
		t.truncateCondition,
	)

	// t.ownerT.Log(sql)

	err := t.conn.Exec(ctx, sql)
	require.NoError(t.ownerT, err)
}

// ApplyMigrations - накатывает миграции расположенные в указанной директории.
func (t *PostgresTester) ApplyMigrations(dirPath string) {
	t.ownerT.Helper()

	pgxPool, err := t.conn.Cli()
	require.NoError(t.ownerT, err)

	db := stdlib.OpenDBFromPool(pgxPool)
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	require.NoError(t.ownerT, err)

	dbMigrate, err := migrate.NewWithDatabaseInstance("file://"+dirPath, postgresDB, driver)
	require.NoError(t.ownerT, err)
	defer dbMigrate.Close()

	err = dbMigrate.Up()
	require.NoError(t.ownerT, err)
}

// ApplyFixtures - загружает данные из указанной директории (имя файла = схема + '.' + имя таблицы) в БД.
// Перед добавлением данных таблица будет очищена.
func (t *PostgresTester) ApplyFixtures(dirPath string) {
	t.ownerT.Helper()

	pgxPool, err := t.conn.Cli()
	require.NoError(t.ownerT, err)

	db := stdlib.OpenDBFromPool(pgxPool)
	defer db.Close()

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(dirPath),
	)
	require.NoError(t.ownerT, err)

	require.NoError(t.ownerT, fixtures.Load())
}

// CountRows - возвращает количество записей указанной таблицы находящейся в текущей схеме.
func (t *PostgresTester) CountRows(ctx context.Context, tableName string) (count int) {
	t.ownerT.Helper()

	err := t.conn.
		QueryRow(ctx, `SELECT COUNT(*) FROM `+tableName).
		Scan(&count)

	require.NoError(t.ownerT, err)

	return count
}

// Destroy - освобождает ресурсы объекта когда он уже больше не нужен.
func (t *PostgresTester) Destroy(ctx context.Context) {
	t.ownerT.Helper()

	require.NoError(t.ownerT, t.container.Terminate(ctx))
}

func newPostgres(ctx context.Context, dsn string) (*mrpostgres.ConnAdapter, error) {
	conn := mrpostgres.New()
	opts := mrpostgres.Options{
		DSN: dsn,
	}

	if err := conn.Connect(ctx, opts); err != nil {
		return nil, err
	}

	return conn, conn.Ping(ctx)
}

func prepareTruncateCondition(dbSchemas, excludedTables []string) (condition string) {
	if len(dbSchemas) == 0 {
		dbSchemas = append(dbSchemas, postgresDefaultSchema)
	}

	condition = " AND relnamespace IN ('" + strings.Join(dbSchemas, "'::regnamespace,'") + "'::regnamespace)"

	if len(excludedTables) > 0 {
		prefix := postgresDefaultSchema + "."

		// публичная схема срезается у всех таблиц, иначе условие работать не будет правильно
		for i := range excludedTables {
			excludedTables[i] = strings.TrimPrefix(excludedTables[i], prefix)
		}

		condition += " AND oid::regclass NOT IN ('" + strings.Join(excludedTables, "'::regclass,'") + "'::regclass)"
	}

	return condition
}
