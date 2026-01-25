package mrredis

import (
	"context"
	"fmt"
	"time"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/redis/go-redis/v9"
)

// go get -u github.com/redis/go-redis/v9

const (
	connectionName = "Redis"
	testKey        = "testKey-d6b6943c-e1b2-4625-b133-9805a5cf5f8d"

	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 5 * time.Second
)

type (
	// ConnAdapter - адаптер для работы с Redis клиентом.
	ConnAdapter struct {
		conn   redis.UniversalClient
		tracer mrtrace.Tracer
	}

	// Options - опции для создания соединения для ConnAdapter.
	Options struct {
		DSN          string // если указано, то Host, Port не используются, но Password более приоритетен если явно указан
		Host         string
		Port         string
		Password     string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
)

// New - создаёт объект ConnAdapter.
func New(tracer mrtrace.Tracer) *ConnAdapter {
	return &ConnAdapter{
		tracer: tracer,
	}
}

// Connect - создаёт соединение с указанными опциями.
func (c *ConnAdapter) Connect(_ context.Context, opts Options) error {
	if c.conn != nil {
		return errors.ErrInternalStorageConnectionIsAlreadyCreated.New("source", connectionName)
	}

	var (
		addr string
		db   int
	)

	if opts.ReadTimeout == 0 {
		opts.ReadTimeout = defaultReadTimeout
	}

	if opts.WriteTimeout == 0 {
		opts.WriteTimeout = defaultWriteTimeout
	}

	if opts.DSN != "" {
		parsedOpts, err := redis.ParseURL(opts.DSN)
		if err != nil {
			return fmt.Errorf("error parsing redis DSN '%s': %w", opts.DSN, err)
		}

		addr = parsedOpts.Addr
		db = parsedOpts.DB

		if opts.Password == "" {
			opts.Password = parsedOpts.Password
		}
	} else {
		addr = fmt.Sprintf("%s:%s", opts.Host, opts.Port)
	}

	c.conn = redis.NewClient(
		&redis.Options{
			Addr:         addr,
			Password:     opts.Password,
			DB:           db,
			ReadTimeout:  opts.ReadTimeout,
			WriteTimeout: opts.WriteTimeout,
		},
	)

	return nil
}

// Ping - сообщает, установлено ли соединение и является ли оно стабильным.
func (c *ConnAdapter) Ping(ctx context.Context) error {
	if c.conn == nil {
		return errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	ping := c.conn.Ping(ctx)

	if err := ping.Err(); err != nil {
		return errors.ErrSystemStorageConnectionFailed.Wrap(err, "source", connectionName)
	}

	if ping.Val() != "PONG" {
		return errors.ErrInternalStorageQueryFailed.WithDetails(
			"unexpected ping response",
			"source", connectionName,
		)
	}

	get := c.conn.Get(ctx, testKey)
	if err := get.Err(); err != nil && !errors.Is(err, redis.Nil) {
		return errors.ErrInternalStorageQueryFailed.Wrap(
			err,
			"source", connectionName,
			"test_key", testKey,
		)
	}

	return nil
}

// Cli - возвращается нативный объект, с которым работает данный адаптер.
func (c *ConnAdapter) Cli() (redis.UniversalClient, error) {
	if c.conn == nil {
		return nil, errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	return c.conn, nil
}

// Close - закрывает текущее соединение.
func (c *ConnAdapter) Close() error {
	if c.conn == nil {
		return errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	if err := c.conn.Close(); err != nil {
		return errors.ErrSystemStorageFailedToClose.Wrap(err, "source", connectionName)
	}

	c.conn = nil

	return nil
}
