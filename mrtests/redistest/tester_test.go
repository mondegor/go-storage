package redistest_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	"github.com/mondegor/go-storage/mrtests/redistest"
)

// TestTester - проверяет, что контейнер Redis поднимается и тестер работает с ним:
// записывает и читает данные и очищает их.
func TestTester(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)

	ctx := context.Background()

	rdt := redistest.NewTester(t)

	defer rdt.Destroy(ctx)

	cli, err := rdt.Conn().Cli()
	require.NoError(t, err)

	require.NoError(t, cli.Set(ctx, "key", "value", time.Minute).Err())

	got, err := cli.Get(ctx, "key").Result()
	require.NoError(t, err)
	assert.Equal(t, "value", got)

	rdt.FlushAll(ctx)

	exists, err := cli.Exists(ctx, "key").Result()
	require.NoError(t, err)
	assert.Zero(t, exists)
}
