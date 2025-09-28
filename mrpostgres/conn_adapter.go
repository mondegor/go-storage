package mrpostgres

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mondegor/go-sysmess/mrerr/mr"

	"github.com/mondegor/go-storage/mrstorage"
)

// go get -u github.com/jackc/pgx/v5

const (
	connectionName = "Postgres"
	driverName     = "postgres"

	defaultMaxConns        = 4
	defaultMaxConnLifetime = time.Hour
	defaultMaxConnIdleTime = 30 * time.Minute
	defaultConnTimeout     = 60 * time.Second
)

type (
	// ConnAdapter - адаптер для работы с Postgres клиентом.
	ConnAdapter struct {
		pool *pgxpool.Pool
	}

	// Options - опции для создания соединения для ConnAdapter.
	Options struct {
		DSN             string // если указано, то Host, Port, Database, Username не используются, но Password более приоритетен если явно указан
		Host            string
		Port            string
		Database        string
		Username        string
		Password        string
		MaxPoolSize     int
		MaxConnLifetime time.Duration
		MaxConnIdleTime time.Duration
		ConnTimeout     time.Duration
		QueryTracer     pgx.QueryTracer
		AfterConnect    func(ctx context.Context, conn *pgx.Conn) error
	}
)

// New - создаёт объект ConnAdapter.
func New() *ConnAdapter {
	return &ConnAdapter{}
}

// Connect - создаёт пул соединений с указанными опциями.
func (c *ConnAdapter) Connect(ctx context.Context, opts Options) error {
	if c.pool != nil {
		return mr.ErrStorageConnectionIsAlreadyCreated.New(connectionName)
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
		return mr.ErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	c.pool = pool

	return nil
}

// Ping - проверяет работоспособность пула соединений.
func (c *ConnAdapter) Ping(ctx context.Context) error {
	if c.pool == nil {
		return mr.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	if err := c.pool.Ping(ctx); err != nil {
		return mr.ErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	var maxValue uint64

	row := c.pool.QueryRow(ctx, `SELECT 18446744073709551615`)
	if err := row.Scan(&maxValue); err != nil {
		return wrapErrorFetchDataFailed(err)
	}

	return nil
}

// HijackConn - извлекает соединение из пула, которое будет использоваться
// независимо от него и должно быть закрыто тем, кто вызвал данный метод.
func (c *ConnAdapter) HijackConn(ctx context.Context) (*pgx.Conn, error) {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, mr.ErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	return conn.Hijack(), nil
}

// Cli - возвращается нативный объект, с которым работает данный адаптер.
func (c *ConnAdapter) Cli() (*pgxpool.Pool, error) {
	if c.pool == nil {
		return nil, mr.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	return c.pool, nil
}

// Close - закрывает пул соединений.
func (c *ConnAdapter) Close() error {
	if c.pool == nil {
		return mr.ErrStorageConnectionIsNotOpened.New(connectionName)
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
