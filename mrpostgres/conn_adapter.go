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

func (c *ConnAdapter) Connect(opt Options) error {
	if c.pool != nil {
		return mrcore.FactoryErrStorageConnectionIsAlreadyCreated.New(connectionName)
	}

	cnf, err := pgxpool.ParseConfig(getConnString(&opt))

	if err != nil {
		return err
	}

	cnf.MaxConns = int32(opt.MaxPoolSize)
	cnf.ConnConfig.ConnectTimeout = opt.ConnTimeout
	cnf.MaxConnLifetime = opt.ConnTimeout

	if opt.AfterConnectFunc != nil {
		pgxFunc, ok := opt.AfterConnectFunc().(pgxConnectFunc)

		if !ok {
			return mrcore.FactoryErrInternalTypeAssertion.New("pgxConnectFunc", opt.AfterConnectFunc())
		}

		cnf.AfterConnect = pgxFunc
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cnf)

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

func getConnString(o *Options) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		o.Username,
		o.Password,
		o.Host,
		o.Port,
		o.Database)
}
