# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

GoStorage is a Go library (`github.com/mondegor/go-storage`, Go 1.26) providing concrete storage adapters for RabbitMQ (`amqp091-go`), Redis (`go-redis/v9`, plus `redislock` for distributed locking), and S3/MinIO + the native filesystem (file providers). It is consumed as a dependency by other projects; it has no `main` (the `examples/` directory holds standalone runnable demos). Doc comments and CHANGELOG are written in Russian — match that language when editing them.

The abstract interfaces these adapters implement (`mrstorage`, `mrlock`, the `mrpostgres` adapter, the `media` model types, logging/tracing) now live in the sibling module **`github.com/mondegor/go-core`**, not in this repo. This repo depends on `go-core`, not vice versa. (A commented-out `replace` directive in `go.mod` points at a local `../go-core` checkout for cross-module development.)

## Commands

The canonical workflow uses the external [`mrcmd`](https://github.com/mondegor/mrcmd) tool, wrapped by `make` targets. `mrcmd` may not be installed locally; the plain `go`/tool equivalents are listed alongside.

| Task | make / mrcmd | Direct equivalent |
|------|--------------|-------------------|
| Run all tests | `make test` | `go test ./...` |
| Coverage report (`test-coverage-full.html`) | `make test-report` | — |
| Lint | `make lint` (gofumpt → goimports → gci → golangci-lint); `mrcmd go-dev lint` runs only golangci-lint | `golangci-lint run` (config: `.golangci.yaml`) |
| Format | (part of `make lint`) | `gofumpt -l -w -extra ./`, then `goimports -l -w -local github.com/mondegor/go-storage .` (mrcmd passes only non-generated `.go` files), then `gci write --skip-generated -s standard -s default -s "prefix(github.com/mondegor/go-storage)" .` (same sections as `.golangci.yaml`) |
| Download deps | `make deps` | `go mod download` |
| Upgrade deps | `make deps-upgrade` | `go get -u ./...` then `go mod tidy` |

- Run a single test: `go test ./mrminio/ -run TestFileProvider -v`
- `make check-and-fix` runs the full pre-commit chain: generate → format → lint → test → plantuml. `make full` adds dependency download first; `make full2` upgrades dependencies first. (`make generate` / `go generate ./...` is a no-op here — this repo currently has no `go:generate` directives or generated mocks; the mocked interfaces and their mocks live in `go-core`.)

### Tests and the `mrtests` harness

The tests in this repo (`mrminio`, `mrredis/locker`, `mrfilestorage`) are only interface-conformance checks (`assert.Implements`) and need no external services — `go test ./...` runs without Docker. The `mrtests/` package provides the reusable harness: `mrtests/helpers` wraps testcontainers containers (`postgres_container.go`, `redis_container.go`, `minio_container.go`; a running Docker daemon is required wherever it is used), `mrtests/infra` provides `*Tester` objects (`PostgresTester`, `RedisTester`, `MinioTester`) that handle migrations (`golang-migrate`) and fixture loading/truncation (`testfixtures`). The Postgres harness is retained here for downstream projects' tests even though the Postgres *adapter* itself now lives in `go-core` — it is the only remaining user of `pgx/v5` in this module.

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
- **`mrfilestorage`** — `file_provider.go` (+ `file_system.go`, `errors.go`, `wrapper_errors.go`) implementing `FileProvider` over the local filesystem; `Ping` writes/reads a sentinel `testFile`.
- **`mrrabbitmq`** — `conn_adapter.go` only: AMQP 0.9.1 connection management (`amqp://User:Password@Host:Port/`).

### Dependencies on `go-core`

Adapters import, in addition to `mrstorage`: `go-core/errors` (error taxonomy), `go-core/mrlog` (logging), `go-core/mrtrace` (tracing), `go-core/mrmodel/media` (file model types), and `go-core/util/*` (`casttype`, `mime`). Keep these consistent with the version pinned in `go.mod`.
