package pgtest

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	containerLogOccurrence = 2
	containerLogTimeout    = 5 * time.Second
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
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(containerLogOccurrence).
				WithStartupTimeout(containerLogTimeout),
		),
	)
	if err != nil {
		return nil, err
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
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
