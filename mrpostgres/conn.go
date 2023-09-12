package mrpostgres

import (
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/mondegor/go-webcore/mrcore"
)

// go get -u github.com/jackc/pgx/v5
// go get -u github.com/Masterminds/squirrel

const (
	connectionName = "postgres"
)

type (
    Connection struct {
        pool *pgxpool.Pool
    }

    Options struct {
        Host string
        Port string
        Database string
        Username string
        Password string
        MaxPoolSize int32
        ConnAttempts int32
        ConnTimeout time.Duration
    }
)

func New() *Connection {
    return &Connection{}
}

func (c *Connection) Connect(opt Options) error {
    if c.pool != nil {
        return mrcore.FactoryErrConnectionIsAlreadyCreated.New(connectionName)
    }

    cnf, err := pgxpool.ParseConfig(getConnString(&opt))

    if err != nil {
        return err
    }

    cnf.MaxConns = opt.MaxPoolSize
    cnf.ConnConfig.ConnectTimeout = opt.ConnTimeout
    cnf.MaxConnLifetime = opt.ConnTimeout

    pool, err := pgxpool.NewWithConfig(context.Background(), cnf)

    if err != nil {
        return mrcore.FactoryErrConnectionFailed.Wrap(err, connectionName)
    }

    c.pool = pool

    return nil
}

func (c *Connection) Ping(ctx context.Context) error {
    if c.pool == nil {
        return mrcore.FactoryErrConnectionIsNotOpened.New(connectionName)
    }

    return c.pool.Ping(ctx)
}

func (c *Connection) Cli() *pgxpool.Pool {
    return c.pool
}

func (c *Connection) Close() error {
    if c.pool == nil {
        return mrcore.FactoryErrConnectionIsNotOpened.New(connectionName)
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
