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
	// ConnAdapter - comment struct.
	ConnAdapter struct {
		conn redis.UniversalClient
	}

	// Options - опции для создания соединения для ConnAdapter.
	Options struct {
		Host        string
		Port        string
		Password    string
		ConnTimeout time.Duration
	}
)

// New - создаёт объект ConnAdapter.
func New() *ConnAdapter {
	return &ConnAdapter{}
}

// Connect - comment method.
func (c *ConnAdapter) Connect(_ context.Context, opts Options) error {
	if c.conn != nil {
		return mrcore.ErrStorageConnectionIsAlreadyCreated.New(connectionName)
	}

	c.conn = redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", opts.Host, opts.Port),
			Password: opts.Password,
		},
	)

	return nil
}

// Ping - comment method.
func (c *ConnAdapter) Ping(ctx context.Context) error {
	if c.conn == nil {
		return mrcore.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	if _, err := c.conn.Ping(ctx).Result(); err != nil {
		return mrcore.ErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	return nil
}

// Cli - comment method.
func (c *ConnAdapter) Cli() redis.UniversalClient { //nolint:ireturn
	return c.conn
}

// Close - comment method.
func (c *ConnAdapter) Close() error {
	if c.conn == nil {
		return mrcore.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	if err := c.conn.Close(); err != nil {
		return mrcore.ErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	c.conn = nil

	return nil
}
