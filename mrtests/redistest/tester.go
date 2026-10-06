package redistest

import (
	"context"
	"testing"

	"github.com/mondegor/go-core/mrtrace"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-storage/mrtests/internal/testenv"
)

const (
	dockerImageEnv = "MRTESTS_REDIS_DOCKER_IMAGE"
	dockerImage    = "redis:7.4.8-alpine3.21"
	password       = "123456"
)

type (
	// Tester - вспомогательный объект для работы с тестовой БД.
	Tester struct {
		ownerT    *testing.T
		container *Container
		conn      *mrredis.ConnAdapter
	}
)

// NewTester - создаёт объект Tester.
// Докер образ берётся из DockerImage, Redis запускается с паролем.
func NewTester(t *testing.T) *Tester {
	t.Helper()

	ctx := context.Background()
	container, err := NewContainer(
		ctx,
		DockerImage(),
		password,
	)
	require.NoError(t, err)

	conn, err := newRedis(ctx, container.DSN())
	require.NoError(t, err)

	return &Tester{
		ownerT:    t,
		container: container,
		conn:      conn,
	}
}

// DockerImage - возвращает докер образ Redis из переменной окружения
// MRTESTS_REDIS_DOCKER_IMAGE, а если она не задана, то берётся значение по умолчанию.
func DockerImage() string {
	return testenv.DockerImage(dockerImageEnv, dockerImage)
}

// Conn - возвращает менеджер текущего соединения с БД.
func (t *Tester) Conn() *mrredis.ConnAdapter {
	t.ownerT.Helper()

	return t.conn
}

// FlushAll - очистка всех данных в Redis.
func (t *Tester) FlushAll(ctx context.Context) {
	t.ownerT.Helper()

	redisCli, err := t.conn.Cli()
	require.NoError(t.ownerT, err)

	cmd := redisCli.FlushAll(ctx)
	require.NoError(t.ownerT, cmd.Err())
}

// Destroy - освобождает ресурсы объекта когда он уже больше не нужен.
func (t *Tester) Destroy(ctx context.Context) {
	t.ownerT.Helper()

	require.NoError(t.ownerT, t.container.Terminate(ctx))
}

func newRedis(ctx context.Context, dsn string) (*mrredis.ConnAdapter, error) {
	conn := mrredis.New(mrtrace.NopTracer())
	opts := mrredis.Options{
		DSN:      dsn,
		Password: password,
	}

	if err := conn.Connect(ctx, opts); err != nil {
		return nil, err
	}

	return conn, conn.Ping(ctx)
}
