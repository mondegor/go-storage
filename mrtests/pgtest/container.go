package pgtest

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	containerPort = "5432/tcp"

	// containerLogOccurrence - сколько раз Postgres пишет в лог о готовности
	// (первый раз при инициализации БД, второй - после основного запуска).
	containerLogOccurrence = 2

	// containerStartupTimeout - время ожидания готовности контейнера
	// (при параллельном запуске пакетов (go test -p) одновременно стартуют
	// несколько контейнеров, и их старт растягивается).
	containerStartupTimeout = 30 * time.Second
)

type (
	// Container - обёртка докер контейнера Postgres.
	Container struct {
		*tcpostgres.PostgresContainer
		dsn string
	}
)

// NewContainer - создаёт объект Container.
func NewContainer(ctx context.Context, dockerImage, database, username, password string) (*Container, error) {
	container, err := tcpostgres.Run(
		ctx,
		dockerImage,
		tcpostgres.WithDatabase(database),
		tcpostgres.WithUsername(username),
		tcpostgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort(containerPort).WithStartupTimeout(containerStartupTimeout),
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(containerLogOccurrence).
				WithStartupTimeout(containerStartupTimeout),
		),
	)
	if err != nil {
		// контейнер, не дождавшийся готовности, удаляется сразу, иначе он продолжит
		// работать до конца прогона и нагружать Docker, провоцируя новые таймауты
		if container != nil {
			_ = container.Terminate(context.WithoutCancel(ctx))
		}

		return nil, err
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(context.WithoutCancel(ctx))

		return nil, err
	}

	return &Container{
		PostgresContainer: container,
		dsn:               dsn,
	}, nil
}

// DSN - возвращает строку соединения с контейнером Postgres.
func (c *Container) DSN() string {
	return c.dsn
}
