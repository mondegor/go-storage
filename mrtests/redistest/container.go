package redistest

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

type (
	// Container - обёртка докер контейнера Redis.
	Container struct {
		*tcredis.RedisContainer
		dsn string
	}
)

// NewContainer - создаёт объект Container.
// Если указан password, то сервер запускается с требованием аутентификации,
// иначе используется команда запуска, заданная в докер образе.
func NewContainer(ctx context.Context, dockerImage, password string) (*Container, error) {
	var opts []testcontainers.ContainerCustomizer

	if password != "" {
		opts = append(opts, testcontainers.WithCmd("redis-server", "--requirepass", password))
	}

	container, err := tcredis.Run(ctx, dockerImage, opts...)
	if err != nil {
		return nil, err
	}

	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	return &Container{
		RedisContainer: container,
		dsn:            dsn,
	}, nil
}

// DSN - возвращает строку соединения с контейнером Redis.
func (c *Container) DSN() string {
	return c.dsn
}
