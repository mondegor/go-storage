package redistest

import (
	"context"
	"errors"
	"testing"

	"github.com/mondegor/go-core/mrtrace"
	"github.com/stretchr/testify/assert"
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
	// Tester - вспомогательный объект для работы с тестовым Redis.
	Tester struct {
		conn *mrredis.ConnAdapter
	}
)

// NewTester - создаёт объект Tester.
// Докер образ берётся из DockerImage, Redis запускается с паролем.
// Соединение и контейнер освобождаются автоматически по завершении теста t,
// а методы тестера принимают t того теста (подтеста), из которого они вызываются.
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
	if err != nil {
		_ = container.Terminate(ctx)
	}

	require.NoError(t, err)

	// ресурсы освобождаются по завершении теста t (после всех его подтестов)
	t.Cleanup(func() {
		assert.NoError(t, errors.Join(conn.Close(), container.Terminate(context.Background())))
	})

	return &Tester{
		conn: conn,
	}
}

// DockerImage - возвращает докер образ Redis из переменной окружения
// MRTESTS_REDIS_DOCKER_IMAGE, а если она не задана, то берётся значение по умолчанию.
func DockerImage() string {
	return testenv.DockerImage(dockerImageEnv, dockerImage)
}

// Conn - возвращает адаптер текущего соединения с Redis.
func (rt *Tester) Conn() *mrredis.ConnAdapter {
	return rt.conn
}

// FlushAll - очищает все данные Redis.
func (rt *Tester) FlushAll(t *testing.T, ctx context.Context) {
	t.Helper()

	redisCli, err := rt.conn.Cli()
	require.NoError(t, err)

	cmd := redisCli.FlushAll(ctx)
	require.NoError(t, cmd.Err())
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

	if err := conn.Ping(ctx); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return conn, nil
}
