package pgtest_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	"github.com/mondegor/go-storage/mrtests/pgtest"
)

// TestTester - проверяет, что контейнер Postgres поднимается и тестер работает с БД:
// накатывает миграции, загружает фикстуры, считает записи и очищает таблицы, кроме исключённых.
func TestTester(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)

	ctx := context.Background()

	migrationsDir := t.TempDir()
	writeFile(
		t, migrationsDir, "1_init.up.sql",
		`CREATE TABLE items (id int PRIMARY KEY, caption text NOT NULL);
		CREATE TABLE settings (id int PRIMARY KEY, caption text NOT NULL);`,
	)
	writeFile(t, migrationsDir, "1_init.down.sql", `DROP TABLE settings; DROP TABLE items;`)

	fixturesDir := t.TempDir()
	writeFile(t, fixturesDir, "items.yml", "- id: 1\n  caption: first\n- id: 2\n  caption: second\n")
	writeFile(t, fixturesDir, "settings.yml", "- id: 1\n  caption: first\n")

	pgt := pgtest.NewTester(t, nil, []string{"public.settings"})

	defer pgt.Destroy(ctx)

	var one int
	require.NoError(t, pgt.ConnManager().Conn(ctx).QueryRow(ctx, `SELECT 1`).Scan(&one))
	assert.Equal(t, 1, one)

	pgt.ApplyMigrations(migrationsDir)
	pgt.ApplyFixtures(fixturesDir)
	assert.Equal(t, 2, pgt.CountRows(ctx, "items"))
	assert.Equal(t, 1, pgt.CountRows(ctx, "settings"))

	pgt.TruncateTables(ctx)
	assert.Equal(t, 0, pgt.CountRows(ctx, "items"))
	assert.Equal(t, 1, pgt.CountRows(ctx, "settings"))
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()

	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
}
