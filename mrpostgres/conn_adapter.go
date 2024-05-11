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
)

type (
	ConnAdapter struct {
		pool *pgxpool.Pool
		dbExecHelper
	}

	Options struct {
		Host             string
		Port             string
		Database         string
		Username         string
		Password         string
		MaxPoolSize      int
		ConnAttempts     int
		ConnTimeout      time.Duration
		AfterConnectFunc func() any
	}

	pgxConnectFunc func(ctx context.Context, conn *pgx.Conn) error
)

func New() *ConnAdapter {
	return &ConnAdapter{}
}

func (c *ConnAdapter) Connect(ctx context.Context, opts Options) error {
	if c.pool != nil {
		return mrcore.FactoryErrStorageConnectionIsAlreadyCreated.New(connectionName)
	}

	cfg, err := pgxpool.ParseConfig(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			opts.Username,
			opts.Password,
			opts.Host,
			opts.Port,
			opts.Database,
		),
	)
	if err != nil {
		return err
	}

	cfg.MaxConns = int32(opts.MaxPoolSize)
	cfg.ConnConfig.ConnectTimeout = opts.ConnTimeout
	cfg.MaxConnLifetime = opts.ConnTimeout

	if opts.AfterConnectFunc != nil {
		pgxFunc, ok := opts.AfterConnectFunc().(pgxConnectFunc)

		if !ok {
			return mrcore.FactoryErrInternalTypeAssertion.New("pgxConnectFunc", opts.AfterConnectFunc())
		}

		cfg.AfterConnect = pgxFunc
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return mrcore.FactoryErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	c.pool = pool

	return nil
}

func (c *ConnAdapter) Ping(ctx context.Context) error {
	if c.pool == nil {
		return mrcore.FactoryErrStorageConnectionIsNotOpened.New(connectionName)
	}

	return c.pool.Ping(ctx)
}

func (c *ConnAdapter) Cli() *pgxpool.Pool {
	return c.pool
}

func (c *ConnAdapter) Close() error {
	if c.pool == nil {
		return mrcore.FactoryErrStorageConnectionIsNotOpened.New(connectionName)
	}

	c.pool.Close()
	c.pool = nil

	return nil
}
