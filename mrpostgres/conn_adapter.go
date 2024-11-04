package mrpostgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mondegor/go-webcore/mrcore"
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
		dbExecHelper
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
		ConnAttempts    int
		ConnTimeout     time.Duration
		QueryTracer     pgx.QueryTracer
		AfterConnect    func(ctx context.Context, conn *pgx.Conn) error
	}
)

// New - создаёт объект ConnAdapter.
func New() *ConnAdapter {
	return &ConnAdapter{}
}

// Connect - создаёт соединение с указанными опциями.
func (c *ConnAdapter) Connect(ctx context.Context, opts Options) error {
	if c.pool != nil {
		return mrcore.ErrStorageConnectionIsAlreadyCreated.New(connectionName)
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
		return mrcore.ErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	c.pool = pool

	return nil
}

// Ping - проверяет работоспособность соединения.
func (c *ConnAdapter) Ping(ctx context.Context) error {
	if c.pool == nil {
		return mrcore.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	if err := c.pool.Ping(ctx); err != nil {
		return mrcore.ErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	var maxValue uint64

	row := c.pool.QueryRow(ctx, `SELECT 18446744073709551615`)
	if err := row.Scan(&maxValue); err != nil {
		return wrapErrorFetchDataFailed(err)
	}

	return nil
}

// Cli - возвращается нативный объект, с которым работает данный адаптер.
func (c *ConnAdapter) Cli() *pgxpool.Pool {
	return c.pool
}

// Close - закрывает текущее соединение.
func (c *ConnAdapter) Close() error {
	if c.pool == nil {
		return mrcore.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	c.pool.Close()
	c.pool = nil

	return nil
}
