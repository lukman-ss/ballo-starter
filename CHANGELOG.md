# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project does not adhere to [Semantic Versioning](https://semver.org/).

## [Unreleased] - 2026-11-04

### Changed

- `internal/app/app.go`: Renamed `ballodb` import to `balldb` for consistency with `ballo-database` package naming.
- `internal/http/handler/task.go`: Updated request body reader to use `req.Raw().Body` instead of `req.Body`.

### Added

- `internal/http/handler/health.go`: New health handler with `Healthz` and `Readyz` endpoints.
- `internal/task/task.go`: New in-memory task store implementation.
- `internal/task/task_test.go`: New unit tests for task store (concurrency, validation, CRUD operations).

### Fixed

- `internal/app/app.go`: Fixed database import alias.
- `internal/http/handler/task.go`: Fixed request body access for HTTP request handling.

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
