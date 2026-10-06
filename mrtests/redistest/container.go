package redistest

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	containerPort = "6379/tcp"

	// containerStartupTimeout - время ожидания готовности контейнера
	// (при параллельном запуске пакетов (go test -p) одновременно стартуют
	// несколько контейнеров, и их старт растягивается).
	containerStartupTimeout = 30 * time.Second
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
	opts := []testcontainers.ContainerCustomizer{
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort(containerPort).WithStartupTimeout(containerStartupTimeout),
			wait.ForLog("* Ready to accept connections").WithStartupTimeout(containerStartupTimeout),
		),
	}

	if password != "" {
		opts = append(opts, testcontainers.WithCmd("redis-server", "--requirepass", password))
	}

	container, err := tcredis.Run(ctx, dockerImage, opts...)
	if err != nil {
		// контейнер, не дождавшийся готовности, удаляется сразу, иначе он продолжит
		// работать до конца прогона и нагружать Docker, провоцируя новые таймауты
		if container != nil {
			_ = container.Terminate(context.WithoutCancel(ctx))
		}

		return nil, err
	}

	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		_ = container.Terminate(context.WithoutCancel(ctx))

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
