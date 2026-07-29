# AGENTS.md

OpenTelemetry-compatible Zap logger for Sourcegraph (`github.com/sourcegraph/log`). This is a Go library consumed by other Go services.

## Setup

Requires Go as specified in `go.mod` (`go 1.25.9`). Dependencies are managed with Go modules.

```bash
go mod download
```

## Build, test, lint

These are the commands used in CI (`.github/workflows/pipeline.yml`):

```bash
go build -v ./...
go test -v -race -covermode=atomic ./...
```

Run a single package's tests, e.g.:

```bash
go test -race ./logtest/...
```

Some tests use `github.com/hexops/autogold/v2` for golden-file assertions. Regenerate golden values with:

```bash
go test ./... -update
```

## Layout

- Root package `log` — the public API: `Logger`, `Init`, fields (`fields.go`), levels (`levels.go`), sinks (`sinks*.go`). Start at `doc.go` and `logger.go`.
- `internal/` — implementation details not part of the public API (`encoders`, `globallogger`, `otelfields`, `sinkcores`, `configurable`, `stderr`).
- `logr/` — adapter for `go-logr/logr`; this is a **separate Go module** with its own `go.mod`.
- `logtest/` — testing utilities for code that uses this library.
- `std/` — `log.Logger`/standard-library adapter.
- `output/`, `hook/` — output formatting and Zap hook helpers.
- `cmd/example/` — runnable usage example.

## Conventions

- No global logger: loggers are passed in from callers to preserve scope and fields. Do not introduce package-level or compile-time loggers (see `doc.go`).
- Scopes should be static, CamelCased values, not dynamic identifiers.
- Behavior is configured via environment variables: `SRC_DEVELOPMENT`, `SRC_LOG_FORMAT` (`json`, `json_gcp`, `condensed`), `SRC_LOG_LEVEL` (`debug`, `info`, `warn`, `error`, `none`), and `SRC_LOG_SCOPE_LEVEL` (see `init.go`).
- When changing the `logr/` adapter, run and update tests within that module separately.
