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
)

type (
	// ConnAdapter - адаптер для работы с Postgres клиентом.
	ConnAdapter struct {
		pool *pgxpool.Pool
		dbExecHelper
	}

	// Options - опции для создания соединения для ConnAdapter.
	Options struct {
		DSN          string // если указано, то Host, Port, Database, Username, Password не используются
		Host         string
		Port         string
		Database     string
		Username     string
		Password     string
		MaxPoolSize  int
		ConnAttempts int
		ConnTimeout  time.Duration
		QueryTracer  pgx.QueryTracer
		AfterConnect func(ctx context.Context, conn *pgx.Conn) error
	}
)

// New - создаёт объект ConnAdapter.
func New() *ConnAdapter {
	return &ConnAdapter{}
}

// Connect - comment method.
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

	cfg, err := pgxpool.ParseConfig(opts.DSN)
	if err != nil {
		return err
	}

	cfg.MaxConns = int32(opts.MaxPoolSize)
	cfg.ConnConfig.ConnectTimeout = opts.ConnTimeout
	cfg.MaxConnLifetime = opts.ConnTimeout
	cfg.ConnConfig.Tracer = opts.QueryTracer
	cfg.AfterConnect = opts.AfterConnect

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return mrcore.ErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	c.pool = pool

	return nil
}

// Ping - comment method.
func (c *ConnAdapter) Ping(ctx context.Context) error {
	if c.pool == nil {
		return mrcore.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	return c.pool.Ping(ctx)
}

// Cli - comment method.
func (c *ConnAdapter) Cli() *pgxpool.Pool {
	return c.pool
}

// Close - comment method.
func (c *ConnAdapter) Close() error {
	if c.pool == nil {
		return mrcore.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	c.pool.Close()
	c.pool = nil

	return nil
}
