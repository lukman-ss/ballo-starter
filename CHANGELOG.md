# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project does not adhere to [Semantic Versioning](https://semver.org/).

## [Unreleased] - 2026-11-04

### Changed

- `internal/http/handler/task.go`: Read request body via `req.Body` (embedded `*http.Request`) instead of `req.Request.Body`.

### Added

- `internal/http/handler/health.go`: New health handler with `Healthz` and `Readyz` endpoints.
- `internal/task/task.go`: New in-memory task store implementation.
- `internal/task/task_test.go`: New unit tests for task store (concurrency, validation, CRUD operations).

### Fixed

- `internal/http/handler/task.go`: Fixed request body access for HTTP request handling.

## [Unreleased]

### Changed

- `internal/config/config.go`: Removed unused database configuration fields (`DATABASE_DRIVER`, `DATABASE_DSN`, `DATABASE_MAX_CONNS`, `DATABASE_MAX_IDLE`).
- `internal/app/app.go`: Removed database initialization comment (database wiring removed).

## [0.1.0] - 2026-01-01

Initial release. Verified implementation with:

- Health probes (`/healthz`, `/readyz`)
- Task CRUD endpoints (`/api/v1/tasks`)
- Configuration via environment variables
- JSON response handling
- Template rendering via `ballo-render`
- Graceful shutdown with signal handling
- HTTP middleware (logger, recovery)
- Comprehensive test suite
