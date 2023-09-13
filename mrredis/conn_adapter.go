package mrredis

import (
    "context"
    "fmt"
    "time"

    "github.com/go-redsync/redsync/v4"
    "github.com/go-redsync/redsync/v4/redis/goredis/v9"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/redis/go-redis/v9"
)

// go get -u github.com/redis/go-redis/v9
// go get github.com/go-redsync/redsync/v4

const (
	connectionName = "redis"
)

type (
    ConnAdapter struct {
        conn redis.UniversalClient
        sync *redsync.Redsync
    }

    Options struct {
        Host string
        Port string
        Password string
        ConnTimeout time.Duration
    }
)

func New() *ConnAdapter {
    return &ConnAdapter{}
}

func (c *ConnAdapter) Connect(opt Options) error {
    if c.conn != nil {
        return mrcore.FactoryErrStorageConnectionIsAlreadyCreated.New(connectionName)
    }

    conn := redis.NewClient(getOptions(&opt))
    _, err := conn.Ping(context.Background()).Result()

    if err != nil {
        return mrcore.FactoryErrStorageConnectionFailed.Wrap(err, connectionName)
    }

    pool := goredis.NewPool(conn)

    c.conn = conn
    c.sync = redsync.New(pool)

    return nil
}

func (c *ConnAdapter) Ping(ctx context.Context) error {
    if c.conn == nil {
        return mrcore.FactoryErrStorageConnectionIsNotOpened.New(connectionName)
    }

    _, err := c.conn.Ping(ctx).Result()

    if err != nil {
        return mrcore.FactoryErrStorageConnectionFailed.Wrap(err, connectionName)
    }

    return nil
}

func (c *ConnAdapter) Cli() redis.UniversalClient {
    return c.conn
}

func (c *ConnAdapter) NewMutex(name string, options ...redsync.Option) *redsync.Mutex {
    return c.sync.NewMutex(name, options...)
}

func (c *ConnAdapter) Close() error {
    if c.conn == nil {
        return mrcore.FactoryErrStorageConnectionIsNotOpened.New(connectionName)
    }

    err := c.conn.Close()

    if err != nil {
        return mrcore.FactoryErrStorageConnectionFailed.Wrap(err, connectionName)
    }

    c.conn = nil

    return nil
}

func getOptions(o *Options) *redis.Options {
    return &redis.Options{
        Addr: fmt.Sprintf("%s:%s", o.Host, o.Port),
        Password: o.Password,
    }
}
