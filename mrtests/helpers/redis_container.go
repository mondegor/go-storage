package helpers

import (
	"context"

	"github.com/testcontainers/testcontainers-go/modules/redis"
)

type (
	// RedisContainer - обёртка докер контейнера Redis.
	RedisContainer struct {
		*redis.RedisContainer
		dsn string
	}
)

// NewRedisContainer - создаёт объект RedisContainer.
func NewRedisContainer(ctx context.Context, dockerImage string) (*RedisContainer, error) {
	container, err := redis.Run(
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

	return &RedisContainer{
		RedisContainer: container,
		dsn:            dsn,
	}, nil
}

// DSN - возвращает строку соединения с контейнером Redis.
func (h *RedisContainer) DSN() string {
	return h.dsn
}
