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
	connectionName = "redis"
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

func (c *ConnAdapter) Connect(opt Options) error {
	if c.conn != nil {
		return mrcore.FactoryErrStorageConnectionIsAlreadyCreated.New(connectionName)
	}

	c.conn = redis.NewClient(getOptions(&opt))

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

func getOptions(o *Options) *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", o.Host, o.Port),
		Password: o.Password,
	}
}
