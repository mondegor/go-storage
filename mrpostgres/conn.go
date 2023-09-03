package mrpostgres

import (
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/mondegor/go-storage/mrstorage"
)

// go get -u github.com/jackc/pgx/v5
// go get -u github.com/Masterminds/squirrel

const ConnectionName = "postgres"

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

func (c *Connection) Connect(ctx context.Context, opt Options) error {
    if c.pool != nil {
        return mrstorage.FactoryConnectionIsAlreadyCreated.New(ConnectionName)
    }

    cnf, err := pgxpool.ParseConfig(getConnString(&opt))

    if err != nil {
        return err
    }

    cnf.MaxConns = opt.MaxPoolSize
    cnf.ConnConfig.ConnectTimeout = opt.ConnTimeout
    cnf.MaxConnLifetime = opt.ConnTimeout

    c.pool, err = pgxpool.NewWithConfig(ctx, cnf)

    if err != nil {
        return mrstorage.FactoryConnectionFailed.Wrap(err, ConnectionName)
    }

    return nil
}

func (c *Connection) Ping(ctx context.Context) error {
    return c.pool.Ping(ctx)
}

func (c *Connection) Close() error {
    if c.pool == nil {
        return mrstorage.FactoryConnectionIsNotOpened.New(ConnectionName)
    }

    c.pool.Close()

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
