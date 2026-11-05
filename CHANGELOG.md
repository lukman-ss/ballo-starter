# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project does not adhere to [Semantic Versioning](https://semver.org/).

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
