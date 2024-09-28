package helpers

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresLogOccurrence = 2
	postgresLogTimeout    = 5 * time.Second
)

type (
	// PostgresContainer - обёртка докер контейнера Postgres.
	PostgresContainer struct {
		*postgres.PostgresContainer
		dsn string
	}
)

// NewPostgresContainer - создаёт объект PostgresContainer.
func NewPostgresContainer(ctx context.Context, dockerImage, database, username, password string) (*PostgresContainer, error) {
	container, err := postgres.Run(
		ctx,
		dockerImage,
		postgres.WithDatabase(database),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(postgresLogOccurrence).
				WithStartupTimeout(postgresLogTimeout),
		),
	)
	if err != nil {
		return nil, err
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: container,
		dsn:               dsn,
	}, nil
}

// DSN - возвращает строку соединения с контейнером Postgres.
func (h *PostgresContainer) DSN() string {
	return h.dsn
}
