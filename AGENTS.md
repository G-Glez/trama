# TRAMA — Tournament Records and Metrics Assistant

## What is TRAMA

TRAMA (**T**ournament **R**ecords **A**nd **M**etrics **A**ssistant) is a Go REST API backend for managing tabletop/wargaming tournament records. It handles game system catalogs (systems, editions, factions).

## Technologies

- **Go 1.26** — language
- **Gin** — HTTP web framework
- **SQLite** (pure-Go driver) — embedded database
- **UUID** — entity ID generation
- **Swaggo** — Swagger/OpenAPI 2.0 docs (generated + served)
- **Docker / Docker Compose** — containerization

## Commands

Run all commands via `make`. See [Makefile](./Makefile).

| Command | Description |
|---|---|---|
| `make build-api` | Build the API binary locally |
| `make build-cli` | Build the CLI binary locally |
| `make test` | Run all tests |
| `make migrate` | Build CLI and run DB migrations |
| `make lint` | Run linter |
| `make swagger` | Regenerate Swagger docs |
| `make sqlc` | Regenerate sqlc DAO code from SQL queries |
| `make run-api` | Build Docker image and start container |
| `make stop` | Stop the container |
| `make clean` | Clean build artifacts |

## API Docs (OpenAPI / Swagger)

- Generated spec: `docs/swagger.yaml` and `docs/swagger.json`
- Swagger UI served at `GET /swagger/*any` (debug mode only)
- Add/update annotations in handler files, then run `make swagger`

## Conventions

- **Language**: Go — standard formatting (`gofumpt`), idiomatic naming
- **Architecture**: domain-driven, hexagonal-style — `internal/core/` for domain, `internal/api/` for delivery
- **Imports**: stdlib first, then third-party, then internal; group with blank lines
- **Errors**: return early, wrap with `fmt.Errorf("context: %w", err)`
- **Testing**: `*_test.go` next to implementation; use `httptest` for HTTP handler tests
- **SQLite migrations**: additive only, SQL files in `internal/migrate/migrations/`, runner in `internal/migrate/migrate.go`
- **No commented-out code** — delete it
- **No hardcoded secrets** — use env vars loaded in `cmd/*/provider/provider.go` via `caarlos0/env`
- **Commits**: [Conventional Commits](https://www.conventionalcommits.org/) — `feat:` for features, `fix:` for bug fixes, `mid:` for intermediate/wip commits

## Configuration

Via environment variables (see `.local.env`):

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | Server listen port |
| `GIN_MODE` | `debug` | Gin mode (`debug`/`release`/`test`) |
| `DATABASE_PATH` | `./data/trama.db` | SQLite database file path |

## Git Policy

**Do not** stage, commit, amend, push, merge, rebase, or create branches/PRs unless the user explicitly grants permission to do so.

## Branching Strategy

- **`main`** — production branch. Protected: requires PR with 1 approval, no direct pushes, no force push, linear history enforced.
- **`dev`** — development/integration branch.
- **`feat/*`** — feature branches. Create from `dev`, PR into `dev` when ready.

Workflow: `feat/lo-que-sea` → PR to `dev` → integration → PR to `main` → production.

