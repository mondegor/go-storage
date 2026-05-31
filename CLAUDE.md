# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

GoStorage is a Go library (`github.com/mondegor/go-storage`, Go 1.25) providing storage adapters for PostgreSQL (`pgx/v5`), RabbitMQ (`amqp091-go`), Redis (`go-redis/v9` + `redislock`), S3/MinIO, and the native filesystem. It is consumed as a dependency by other projects; it has no `main` (the `examples/` directory holds standalone runnable demos). Doc comments and CHANGELOG are written in Russian — match that language when editing them.

## Commands

The canonical workflow uses the external [`mrcmd`](https://github.com/mondegor/mrcmd) tool, wrapped by `make` targets. `mrcmd` may not be installed locally; the plain `go`/tool equivalents are listed alongside.

| Task | make / mrcmd | Direct equivalent |
|------|--------------|-------------------|
| Run all tests | `make test` | `go test ./...` |
| Coverage report (`test-coverage-full.html`) | `make test-report` | — |
| Lint | `make lint` | `golangci-lint run` (config: `.golangci.yaml`) |
| Format | (part of `make lint`) | `gofumpt -l -w -extra ./` then `goimports -d -local github.com/mondegor/go-storage ./` |
| Regenerate mocks | `make generate` | `go generate ./...` |
| Download deps | `make deps` | `go mod download` |

- Run a single test: `go test ./mrpostgres/builder/part/ -run TestSQLCondition -v`
- `make check-and-fix` runs the full pre-commit chain: generate → format → lint → test → plantuml. `make full` adds dependency download first.

### Tests require Docker

Integration tests (e.g. `mrpostgres`, `mrminio`, `mrlock/redislocker`, `mrfilestorage`) spin up real services via `testcontainers-go`, so a running Docker daemon is required. There are **no build tags** separating unit from integration tests — `go test ./...` will attempt to start containers. The `mrtests/` package provides the reusable harness: `mrtests/helpers` wraps testcontainers containers, `mrtests/infra` provides `*Tester` objects (e.g. `PostgresTester`) that handle migrations (`golang-migrate`) and fixture loading/truncation (`testfixtures`).

## Architecture

### `mrstorage` is the interface hub

The root `mrstorage` package defines the abstract interfaces that everything else implements or consumes — it imports no concrete drivers. Adapter packages depend on `mrstorage`, not vice versa. Key interfaces:

- **`DBConnManager`** (`db.go`) — combines `DBConn` (Query/QueryRow/Exec) with `DBTxManager.Do(...)`. This is the primary type that repository code depends on.
- **`SQLBuilder`** (`sql_builder.go`) — fluent builder producing `SQLPart`s for SET / WHERE+JOIN / ORDER BY / LIMIT. The central concept is `SQLPartFunc func(argumentNumber int) (sql string, args []any)`: parts are evaluated lazily so that placeholder numbering (`$1, $2, …`) stays correct when parts are concatenated.
- **`FileProvider`** (`file_provider.go`) — uniform file API (Info/Download/Upload/Remove) implemented by both `mrminio` (S3) and `mrfilestorage` (local FS). `FileProviderPool` registers multiple providers by name.
- **`SequenceGenerator`**, **`DBStatProvider`** — ID generation and pool metrics.

### Transaction propagation via context

`mrpostgres.ConnManager` stores the active `pgx.Tx` in the `context.Context` under a private key. `Conn(ctx)` returns the transaction if one is in flight, otherwise the raw pool connection — so callers write the same code inside and outside a transaction. `Do()` detects nesting: a nested `Do` reuses the existing transaction and only logs a warning if it requests a *higher* isolation level than the outer one. Isolation levels are an internal enum (`mrstorage/txisolevel`) mapped to pgx options, set via `WithTxIsoLevel*` options (default `ReadCommitted`).

### Adapter package conventions

Each driver package (`mrpostgres`, `mrredis`, `mrminio`, `mrrabbitmq`, `mrfilestorage`) follows the same shape:
- `conn_adapter.go` — connection lifecycle (Connect/Ping/Close).
- `wrapper_errors.go` — translates driver-specific errors into the shared `go-sysmess/errors` taxonomy (e.g. Postgres `23505` → `ErrInternalStorageDuplicateKeyViolation`). When adding error handling, wrap through these helpers rather than returning raw driver errors.

`mrpostgres` has notable sub-packages: `builder/` (the `SQLBuilder` implementation, split into `part/` low-level and `helper/` high-level), `db/` (generic typed query helpers like `FieldFetcher[RowID, FieldValue]`, `RowExistChecker`, `RowSoftDeleter`), `sequence/`, `listennotify/` (LISTEN/NOTIFY support), and `monitoring/` (query tracer).

### Other packages

- `mrlock` — distributed lock `Locker` interface with `redislocker`, `mutexlocker`, `noplocker` implementations and a generated `mock/`.
- `mrsql` — reflection-based entity metadata: `ParseEntity` reads struct tags `db`, `upd`, `sort` to drive dynamic UPDATE SET and ORDER BY generation (`EntityMeta`).
- `mrstorage/mrprometheus` — Prometheus collector wrapping `DBStatProvider`.
- `mrentity` — shared entity value types.

### Mock generation

Mocks use `golang/mock` (mockgen v1.6.0). Sources carry `//go:generate mockgen -source=... -destination=./mock/...` directives; run `make generate` after changing a mocked interface.
