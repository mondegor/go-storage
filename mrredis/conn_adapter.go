package mrredis

import (
	"context"
	"fmt"
	"time"

	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mrtrace"
	"github.com/redis/go-redis/v9"
)

// go get -u github.com/redis/go-redis/v9

const (
	// connectionName - имя подключения для логирования и трассировки.
	connectionName = "Redis"

	// testKey - ключ для проверки работоспособности соединения в методе Ping.
	testKey = "testKey-d6b6943c-e1b2-4625-b133-9805a5cf5f8d"

	// defaultReadTimeout - таймаут чтения из Redis по умолчанию.
	defaultReadTimeout = 5 * time.Second

	// defaultWriteTimeout - таймаут записи в Redis по умолчанию.
	defaultWriteTimeout = 5 * time.Second
)

type (
	// ConnAdapter - адаптер для работы с Redis клиентом.
	// Предоставляет методы для подключения, выполнения команд (GET, SET, DELETE),
	// проверки работоспособности и получения нативного клиента.
	ConnAdapter struct {
		conn   redis.UniversalClient
		tracer mrtrace.Tracer
	}

	// Options - опции для создания соединения в ConnAdapter.
	// Позволяет подключаться либо по DSN, либо по отдельным параметрам Host/Port.
	Options struct {
		DSN          string        // DSN - строка подключения (если указана, Host и Port игнорируются)
		Host         string        // Host - адрес сервера Redis (используется, если DSN не указан)
		Port         string        // Port - порт сервера Redis (используется, если DSN не указан)
		Password     string        // Password - пароль для аутентификации (переопределяет пароль из DSN, если указан)
		ReadTimeout  time.Duration // ReadTimeout - таймаут чтения (по умолчанию: 5 секунд)
		WriteTimeout time.Duration // WriteTimeout - таймаут записи (по умолчанию: 5 секунд)
	}
)

// New - создаёт объект ConnAdapter без активного соединения.
func New(tracer mrtrace.Tracer) *ConnAdapter {
	return &ConnAdapter{
		tracer: tracer,
	}
}

// Connect - устанавливает соединение с сервером Redis по указанным опциям.
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

// Ping - проверяет работоспособность соединения с Redis.
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

// Cli - возвращает нативный клиент Redis (redis.UniversalClient) для прямого доступа к API.
func (c *ConnAdapter) Cli() (redis.UniversalClient, error) {
	if c.conn == nil {
		return nil, errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	return c.conn, nil
}

// Close - закрывает соединение с Redis и очищает ссылку на клиент.
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
