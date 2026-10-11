# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

GoStorage is a Go library (`github.com/mondegor/go-storage`, Go 1.26) providing concrete storage adapters for RabbitMQ (`amqp091-go`), Redis (`go-redis/v9`, plus `redislock` for distributed locking), and S3/MinIO + the native filesystem (file providers). It is consumed as a dependency by other projects; it has no `main` (the `examples/` directory holds standalone runnable demos). Doc comments and CHANGELOG are written in Russian — match that language when editing them.

The abstract interfaces these adapters implement (`mrstorage`, `mrlock`, the `mrpostgres` adapter, the `media` model types, logging/tracing) now live in the sibling module **`github.com/mondegor/go-core`**, not in this repo. This repo depends on `go-core`, not vice versa. (A commented-out `replace` directive in `go.mod` points at a local `../go-core` checkout for cross-module development.)

## Commands

Tooling is driven by [`mrcmd`](https://github.com/mondegor/mrcmd) with the `go-dev` plugin; `make` targets wrap it. When `mrcmd` is not installed, use the direct equivalents.

| Task | make | mrcmd | Without mrcmd |
|------|------|-------|---------------|
| Download deps | `make deps` | `mrcmd go-dev deps` | `go get ./... && go mod tidy` |
| Upgrade deps | `make deps-upgrade` | `mrcmd go-dev get -u ./...` + `mrcmd go-dev tidy` | `go get -u ./... && go mod tidy` |
| Codegen | `make generate` | `mrcmd go-dev generate` | `go generate ./...` |
| Format + lint | `make lint` | `mrcmd go-dev fmt` + `mrcmd go-dev lint` | `golangci-lint fmt && golangci-lint run ./...` |
| Format diff only | — | `mrcmd go-dev fmt-diff` | `golangci-lint fmt -d` |
| Validate lint config | `make lint-verify` | `mrcmd go-dev lint-verify` | `golangci-lint config verify` (needs network access to golangci-lint.run; not part of CI) |
| Run all tests | `make test` | `mrcmd go-dev test` | `go test -cover ./...` |
| Coverage report (`test-coverage-full.html`) | `make test-report` | `mrcmd go-dev test-report` | — |

- Extra local make targets can be added in a gitignored `Makefile.mk` (included by `Makefile` if present).
- Formatters (gofmt, goimports, gci, gofumpt) and linters are both configured in `.golangci.yaml`; `golangci-lint fmt` skips generated files.
- Dev tools are installed via `mrcmd go-dev install-tools`: `golangci-lint` is pinned in `.env.dist` to the version CI uses (bump both together); `mockgen` and `gotext` are disabled in go-dev and not used by this library.
- `make generate` / `go generate ./...` is a no-op here — this repo currently has no `go:generate` directives or generated mocks; the mocked interfaces and their mocks live in `go-core`.
- Run a single test: `go test ./mrminio/ -run TestFileProvider -v`

### Tests and the `mrtests` harness

The tests in `mrminio`, `mrredis/locker` and `mrfilestorage` are only interface-conformance checks (`assert.Implements`) and need no external services. The `mrtests/` harness is split per service — `mrtests/pgtest` and `mrtests/redistest` (the `test` suffix avoids clashing with `go-redis`/driver package names in consumers). Each has `container.go` (`Container` wrapping the testcontainers module — imported as `tcpostgres`/`tcredis` — plus `NewContainer`, `DSN()`), `tester.go` (`Tester`, `NewTester`, `DockerImage()`) and `tester_test.go` (an integration test against a real container; never skipped — it requires a running Docker and fails without it; CI on `ubuntu-latest` provides Docker). `pgtest.Tester` handles migrations (`golang-migrate`) and fixture loading/truncation (`testfixtures`). Testers default to public Docker Hub images (`postgres:18.3-alpine3.23`, `redis:7.4.8-alpine3.21`) so they run in CI; container setup (e.g. Redis `--requirepass` via the `password` arg of `redistest.NewContainer`) is done in code, never baked into an image. Images are overridable via `MRTESTS_POSTGRES_DOCKER_IMAGE` / `MRTESTS_REDIS_DOCKER_IMAGE`; `pgtest.DockerImage()` / `redistest.DockerImage()` resolve them via `mrtests/internal/testenv.DockerImage`. There is no MinIO harness. The Postgres harness is retained here for downstream projects' tests even though the Postgres *adapter* itself now lives in `go-core` — it is the only remaining user of `pgx/v5` in this module.

## Architecture

### Interfaces live in `go-core`, implementations live here

The abstract interface hub is `github.com/mondegor/go-core/mrstorage`. The adapters in this repo implement those interfaces structurally (no explicit `var _ mrstorage.X` assertions in most files). Relevant external interfaces:

- **`FileProvider`** — uniform file API (Info/Download/Upload/Remove/Ping) implemented by both `mrminio` (S3-compatible MinIO) and `mrfilestorage` (local FS). File metadata flows through `go-core/mrmodel/media` types.
- **`mrlock.Locker`** (in `go-core/mrlock`) — distributed lock interface implemented here by `mrredis/locker`.

### Adapter package conventions

Driver packages (`mrredis`, `mrminio`, `mrrabbitmq`, `mrfilestorage`) share a common shape:
- `conn_adapter.go` (`mrredis`, `mrminio`, `mrrabbitmq`) — connection lifecycle (Connect/Ping/Close) and an `Options` struct. Each adapter logs/traces under a `connectionName`/`providerName` constant.
- `wrapper_errors.go` (`mrredis`, `mrredis/locker`, `mrminio`, `mrfilestorage`) — translates driver-specific errors into the shared `go-core/errors` taxonomy. When adding error handling, wrap through these helpers rather than returning raw driver errors.

Package specifics:
- **`mrredis`** — `conn_adapter.go` (lifecycle, default read/write timeouts), `conn_cmd.go` (command helpers), and `locker/` — the distributed `Locker` built on `bsm/redislock` (`locker.Adapter`, constructed via `NewAdapter(conn, logger, tracer)`), translating `redislock` errors through its own `wrapper_errors.go`.
- **`mrminio`** — `file_provider.go` implementing `FileProvider` over `minio-go/v7`.
- **`mrfilestorage`** — `file_provider.go` (+ `file_system.go`, `errors.go`, `wrapper_errors.go`) implementing `FileProvider` over the local filesystem; `Ping` creates and removes a sentinel `testFile`.
- **`mrrabbitmq`** — `conn_adapter.go` only: AMQP 0.9.1 connection management (`amqp://User:Password@Host:Port/`).

### Dependencies on `go-core`

Adapters import, in addition to `mrstorage`: `go-core/errors` (error taxonomy), `go-core/mrlog` (logging), `go-core/mrtrace` (tracing), `go-core/mrmodel/media` (file model types), and `go-core/util/*` (`casttype`, `mime`). Keep these consistent with the version pinned in `go.mod`.
