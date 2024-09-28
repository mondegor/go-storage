package infra

import (
	"context"
	"testing"

	"github.com/mondegor/go-storage/mrtests/helpers"

	"github.com/mondegor/go-webcore/mrlib"

	"github.com/stretchr/testify/require"

	"github.com/mondegor/go-storage/mrminio"
)

const (
	minioDockerImage = "p/minio:2024-07-13"
	minioUser        = "admin_test"
	minioPassword    = "12345678_test"
)

type (
	// MinioTester - вспомогательный объект для работы с тестовой БД.
	MinioTester struct {
		ownerT      *testing.T
		container   *helpers.MinioContainer
		connAdapter *mrminio.ConnAdapter
	}
)

// NewMinioTester - создаёт объект MinioTester.
func NewMinioTester(t *testing.T, mimeTypes *mrlib.MimeTypeList) *MinioTester {
	t.Helper()

	ctx := context.Background()
	container, err := helpers.NewMinioContainer(
		ctx,
		minioDockerImage,
		minioUser,
		minioPassword,
	)
	require.NoError(t, err)

	connAdapter, err := newMinio(ctx, container.DSN(), mimeTypes)
	require.NoError(t, err)

	return &MinioTester{
		ownerT:      t,
		container:   container,
		connAdapter: connAdapter,
	}
}

// Conn - возвращает менеджер текущего соединения с БД.
func (t *MinioTester) Conn() *mrminio.ConnAdapter {
	t.ownerT.Helper()

	return t.connAdapter
}

// Destroy - освобождает ресурсы объекта когда он уже больше не нужен.
func (t *MinioTester) Destroy(ctx context.Context) {
	t.ownerT.Helper()

	require.NoError(t.ownerT, t.container.Terminate(ctx))
}

func newMinio(ctx context.Context, dsn string, mimeTypes *mrlib.MimeTypeList) (*mrminio.ConnAdapter, error) {
	conn := mrminio.New(false, mimeTypes)
	opts := mrminio.Options{
		DSN:      dsn,
		User:     minioUser,
		Password: minioPassword,
	}

	if err := conn.Connect(ctx, opts); err != nil {
		return nil, err
	}

	return conn, conn.Ping(ctx)
}
