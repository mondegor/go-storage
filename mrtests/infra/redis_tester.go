package infra

import (
	"context"
	"testing"

	"github.com/mondegor/go-sysmess/mrtrace/noptracer"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-storage/mrtests/helpers"
)

const (
	redisDockerImage = "p/redis:7.2.5"
	redisPassword    = "123456"
)

type (
	// RedisTester - вспомогательный объект для работы с тестовой БД.
	RedisTester struct {
		ownerT    *testing.T
		container *helpers.RedisContainer
		conn      *mrredis.ConnAdapter
	}
)

// NewRedisTester - создаёт объект RedisTester.
func NewRedisTester(t *testing.T) *RedisTester {
	t.Helper()

	ctx := context.Background()
	container, err := helpers.NewRedisContainer(
		ctx,
		redisDockerImage,
	)
	require.NoError(t, err)

	conn, err := newRedis(ctx, container.DSN())
	require.NoError(t, err)

	return &RedisTester{
		ownerT:    t,
		container: container,
		conn:      conn,
	}
}

// Conn - возвращает менеджер текущего соединения с БД.
func (t *RedisTester) Conn() *mrredis.ConnAdapter {
	t.ownerT.Helper()

	return t.conn
}

// FlushAll - очистка всех данных в RedisTester.
func (t *RedisTester) FlushAll(ctx context.Context) {
	t.ownerT.Helper()

	redisCli, err := t.conn.Cli()
	require.NoError(t.ownerT, err)

	cmd := redisCli.FlushAll(ctx)
	require.NoError(t.ownerT, cmd.Err())
}

// Destroy - освобождает ресурсы объекта когда он уже больше не нужен.
func (t *RedisTester) Destroy(ctx context.Context) {
	t.ownerT.Helper()

	require.NoError(t.ownerT, t.container.Terminate(ctx))
}

func newRedis(ctx context.Context, dsn string) (*mrredis.ConnAdapter, error) {
	conn := mrredis.New(noptracer.New())
	opts := mrredis.Options{
		DSN:      dsn,
		Password: redisPassword,
	}

	if err := conn.Connect(ctx, opts); err != nil {
		return nil, err
	}

	return conn, conn.Ping(ctx)
}
