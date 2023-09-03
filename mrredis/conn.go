package mrredis

import (
    "context"
    "fmt"
    "time"

    "github.com/mondegor/go-storage/mrstorage"

    "github.com/go-redsync/redsync/v4"
    "github.com/go-redsync/redsync/v4/redis/goredis/v9"
    redislib "github.com/redis/go-redis/v9"
)

// go get -u github.com/redis/go-redis/v9
// go get github.com/go-redsync/redsync/v4

const ConnectionName = "redis"

type (
    Connection struct {
        conn redislib.UniversalClient
        redsync *redsync.Redsync
    }

    Options struct {
        Host string
        Port string
        Password string
        ConnTimeout time.Duration
    }
)

func New() *Connection {
    return &Connection{}
}

func (c *Connection) Cli() redislib.UniversalClient {
    return c.conn
}

func (c *Connection) Connect(opt Options) error {
    if c.conn != nil {
        return mrstorage.FactoryConnectionIsAlreadyCreated.New(ConnectionName)
    }

    conn := redislib.NewClient(getOptions(&opt))
    _, err := conn.Ping(context.Background()).Result()

    if err != nil {
        return mrstorage.FactoryConnectionFailed.Wrap(err, ConnectionName)
    }

    c.conn = conn

    pool := goredis.NewPool(conn)
    c.redsync = redsync.New(pool)

    return nil
}

func (c *Connection) Close() error {
    if c.conn == nil {
        return mrstorage.FactoryConnectionIsNotOpened.New(ConnectionName)
    }

    conn := c.conn
    c.conn = nil
    err := conn.Close()

    if err != nil {
        return mrstorage.FactoryConnectionFailed.Wrap(err, ConnectionName)
    }

    return nil
}

func (c *Connection) NewMutex(name string, options ...redsync.Option) *redsync.Mutex {
    return c.redsync.NewMutex(name, options...)
}

func getOptions(o *Options) *redislib.Options {
    return &redislib.Options{
        Addr: fmt.Sprintf("%s:%s", o.Host, o.Port),
        Password: o.Password,
    }
}
