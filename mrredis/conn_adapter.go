package mrredis

import (
	"context"
	"fmt"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/redis/go-redis/v9"
)

// go get -u github.com/redis/go-redis/v9

const (
	connectionName = "Redis"
)

type (
	ConnAdapter struct {
		conn redis.UniversalClient
	}

	Options struct {
		Host        string
		Port        string
		Password    string
		ConnTimeout time.Duration
	}
)

func New() *ConnAdapter {
	return &ConnAdapter{}
}

func (c *ConnAdapter) Connect(ctx context.Context, opts Options) error {
	if c.conn != nil {
		return mrcore.FactoryErrStorageConnectionIsAlreadyCreated.New(connectionName)
	}

	c.conn = redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", opts.Host, opts.Port),
			Password: opts.Password,
		},
	)

	return nil
}

func (c *ConnAdapter) Ping(ctx context.Context) error {
	if c.conn == nil {
		return mrcore.FactoryErrStorageConnectionIsNotOpened.New(connectionName)
	}

	if _, err := c.conn.Ping(ctx).Result(); err != nil {
		return mrcore.FactoryErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	return nil
}

func (c *ConnAdapter) Cli() redis.UniversalClient {
	return c.conn
}

func (c *ConnAdapter) Close() error {
	if c.conn == nil {
		return mrcore.FactoryErrStorageConnectionIsNotOpened.New(connectionName)
	}

	if err := c.conn.Close(); err != nil {
		return mrcore.FactoryErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	c.conn = nil

	return nil
}
