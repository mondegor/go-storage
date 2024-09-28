package helpers

import (
	"context"

	"github.com/testcontainers/testcontainers-go/modules/minio"
)

type (
	// MinioContainer - обёртка докер контейнера Minio.
	MinioContainer struct {
		*minio.MinioContainer
		dsn string
	}
)

// NewMinioContainer - создаёт объект MinioContainer.
func NewMinioContainer(ctx context.Context, dockerImage, username, password string) (*MinioContainer, error) {
	container, err := minio.Run(
		ctx,
		dockerImage,
		minio.WithUsername(username),
		minio.WithPassword(password),
	)
	if err != nil {
		return nil, err
	}

	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	return &MinioContainer{
		MinioContainer: container,
		dsn:            dsn,
	}, nil
}

// DSN - возвращает строку соединения с контейнером Minio.
func (h *MinioContainer) DSN() string {
	return h.dsn
}
