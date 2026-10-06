package redistest

import (
	"context"

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
func NewContainer(ctx context.Context, dockerImage string) (*Container, error) {
	container, err := tcredis.Run(
		ctx,
		dockerImage,
	)
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
