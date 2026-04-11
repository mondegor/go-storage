package mrpostgres

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mondegor/go-sysmess/errors"

	"github.com/mondegor/go-storage/mrstorage"
)

// go get -u github.com/jackc/pgx/v5

const (
	// connectionName - имя подключения для логирования и трассировки.
	connectionName = "Postgres"

	// driverName - имя драйвера для формирования строки подключения.
	driverName = "postgres"

	// defaultMaxConns - максимальное количество соединений по умолчанию.
	defaultMaxConns = 4

	// defaultMaxConnLifetime - максимальное время жизни соединения по умолчанию.
	defaultMaxConnLifetime = time.Hour

	// defaultMaxConnIdleTime - максимальное время простоя соединения по умолчанию.
	defaultMaxConnIdleTime = 30 * time.Minute

	// defaultConnTimeout - таймаут подключения по умолчанию.
	defaultConnTimeout = 60 * time.Second
)

type (
	// ConnAdapter - адаптер для работы с PostgreSQL через библиотеку pgx.
	// Предоставляет методы для подключения, выполнения запросов и управления пулом соединений.
	ConnAdapter struct {
		pool *pgxpool.Pool
	}

	// Options - опции для создания соединения в ConnAdapter.
	// Позволяет подключаться либо по DSN, либо по отдельным параметрам.
	Options struct {
		DSN             string                                          // DSN - строка подключения (если указана, Host, Port, Database, Username игнорируются)
		Host            string                                          // Host - адрес сервера PostgreSQL
		Port            string                                          // Port - порт сервера PostgreSQL
		Database        string                                          // Database - имя БД
		Username        string                                          // Username - имя пользователя для аутентификации
		Password        string                                          // Password - пароль для аутентификации (переопределяет пароль из DSN, если указан)
		MaxPoolSize     int                                             // MaxPoolSize - максимальный размер пула соединений (по умолчанию: 4)
		MaxConnLifetime time.Duration                                   // MaxConnLifetime - максимальное время жизни соединения (по умолчанию: 1 час)
		MaxConnIdleTime time.Duration                                   // MaxConnIdleTime - максимальное время простоя соединения (по умолчанию: 30 минут)
		ConnTimeout     time.Duration                                   // ConnTimeout - таймаут подключения (по умолчанию: 60 секунд)
		QueryTracer     pgx.QueryTracer                                 // QueryTracer - трассировщик SQL-запросов
		AfterConnect    func(ctx context.Context, conn *pgx.Conn) error // AfterConnect - функция, вызываемая после каждого нового соединения
	}
)

// New - создаёт объект ConnAdapter без активного пула соединений.
func New() *ConnAdapter {
	return &ConnAdapter{}
}

// Connect - создаёт пул соединений PostgreSQL по указанным опциям.
func (c *ConnAdapter) Connect(ctx context.Context, opts Options) error {
	if c.pool != nil {
		return errors.ErrInternalStorageConnectionIsAlreadyCreated.New("source", connectionName)
	}

	if opts.DSN == "" {
		opts.DSN = fmt.Sprintf(
			"%s://%s:%s@%s:%s/%s",
			driverName,
			opts.Username,
			opts.Password,
			opts.Host,
			opts.Port,
			opts.Database,
		)
	}

	if opts.MaxPoolSize == 0 {
		opts.MaxPoolSize = defaultMaxConns
	}

	if opts.MaxConnLifetime == 0 {
		opts.MaxConnLifetime = defaultMaxConnLifetime
	}

	if opts.MaxConnIdleTime == 0 {
		opts.MaxConnIdleTime = defaultMaxConnIdleTime
	}

	if opts.ConnTimeout == 0 {
		opts.ConnTimeout = defaultConnTimeout
	}

	cfg, err := pgxpool.ParseConfig(opts.DSN)
	if err != nil {
		return err
	}

	if opts.MaxPoolSize < 1 || opts.MaxPoolSize > math.MaxInt32 {
		return errors.New("max pool size is incorrect")
	}

	cfg.MaxConns = int32(opts.MaxPoolSize)
	cfg.MaxConnLifetime = opts.MaxConnLifetime
	cfg.MaxConnIdleTime = opts.MaxConnIdleTime
	cfg.ConnConfig.ConnectTimeout = opts.ConnTimeout
	cfg.ConnConfig.Tracer = opts.QueryTracer
	cfg.AfterConnect = opts.AfterConnect

	if opts.Password != "" {
		cfg.ConnConfig.Password = opts.Password
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return errors.ErrSystemStorageConnectionFailed.Wrap(err, "source", connectionName)
	}

	c.pool = pool

	return nil
}

// Ping - проверяет работоспособность пула соединений PostgreSQL.
func (c *ConnAdapter) Ping(ctx context.Context) error {
	if c.pool == nil {
		return errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	if err := c.pool.Ping(ctx); err != nil {
		return errors.ErrSystemStorageConnectionFailed.Wrap(err, "source", connectionName)
	}

	var maxValue uint64

	row := c.pool.QueryRow(ctx, `SELECT 18446744073709551615`)
	if err := row.Scan(&maxValue); err != nil {
		return wrapErrorFetchDataFailed(err)
	}

	return nil
}

// HijackConn - извлекает соединение из пула для независимого использования.
// Возвращённое соединение не управляется пулом и должно быть закрыто вызывающей стороной.
// Полезно для длительного использования соединения (например: LISTEN/NOTIFY).
func (c *ConnAdapter) HijackConn(ctx context.Context) (*pgx.Conn, error) {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, errors.ErrSystemStorageConnectionFailed.Wrap(err, "source", connectionName)
	}

	return conn.Hijack(), nil
}

// Cli - возвращает нативный пул соединений pgx для прямого доступа к API.
func (c *ConnAdapter) Cli() (*pgxpool.Pool, error) {
	if c.pool == nil {
		return nil, errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	return c.pool, nil
}

// Close - закрывает пул соединений.
func (c *ConnAdapter) Close() error {
	if c.pool == nil {
		return errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	c.pool.Close()
	c.pool = nil

	return nil
}

// Query - отправляет SQL запрос к БД и возвращает результат в виде списка записей.
func (c *ConnAdapter) Query(ctx context.Context, sql string, args ...any) (mrstorage.DBQueryRows, error) {
	rows, err := c.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, wrapError(err)
	}

	return &queryRows{
		rows: rows,
	}, nil
}

// QueryRow - отправляет SQL запрос к БД и возвращает результат в виде одной записи.
func (c *ConnAdapter) QueryRow(ctx context.Context, sql string, args ...any) mrstorage.DBQueryRow {
	return &queryRow{
		row: c.pool.QueryRow(ctx, sql, args...),
	}
}

// Exec - отправляет SQL запрос к БД и исполняет его.
func (c *ConnAdapter) Exec(ctx context.Context, sql string, args ...any) error {
	return wrapErrorCommandTag(c.pool.Exec(ctx, sql, args...))
}
