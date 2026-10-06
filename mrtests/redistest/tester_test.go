package redistest_test

import (
	"context"
	"testing"
	"time"

	"github.com/mondegor/go-core/mrtrace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-storage/mrtests/redistest"
)

// TestTester - проверяет, что контейнер Redis поднимается с паролем и тестер работает с ним:
// записывает и читает данные и очищает их.
func TestTester(t *testing.T) {
	ctx := context.Background()

	rdt := redistest.NewTester(t)

	cli, err := rdt.Conn().Cli()
	require.NoError(t, err)

	require.NoError(t, cli.Set(ctx, "key", "value", time.Minute).Err())

	got, err := cli.Get(ctx, "key").Result()
	require.NoError(t, err)
	assert.Equal(t, "value", got)

	rdt.FlushAll(t, ctx)

	exists, err := cli.Exists(ctx, "key").Result()
	require.NoError(t, err)
	assert.Zero(t, exists)
}

// TestNewContainer_RequiresPassword - проверяет, что контейнер, запущенный с паролем,
// отвергает подключение без пароля именно из-за отсутствия аутентификации.
func TestNewContainer_RequiresPassword(t *testing.T) {
	ctx := context.Background()

	container, err := redistest.NewContainer(ctx, redistest.DockerImage(), "secret_test")
	require.NoError(t, err)

	defer func() { assert.NoError(t, container.Terminate(ctx)) }()

	conn := mrredis.New(mrtrace.NopTracer())
	require.NoError(t, conn.Connect(ctx, mrredis.Options{DSN: container.DSN()}))

	defer func() { _ = conn.Close() }()

	cli, err := conn.Cli()
	require.NoError(t, err)

	assert.ErrorContains(t, cli.Ping(ctx).Err(), "NOAUTH")
}
