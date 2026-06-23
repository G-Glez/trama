# TRAMA — Tournament Records and Metrics Assistant

## What is TRAMA

TRAMA (**T**ournament **R**ecords **A**nd **M**etrics **A**ssistant) is a Go REST API backend deployed as an AWS Lambda behind API Gateway HTTP API. It manages tabletop/wargaming tournament records (systems, editions, factions) and user auth.

## Technologies

- **Go 1.25** — language
- **Gin** — HTTP web framework (runs inside Lambda via `aws-lambda-go-api-proxy`)
- **SQLite** (`modernc.org/sqlite`) — embedded database, stored on Lambda ephemeral `/tmp`
- **Terraform** — infrastructure as code (S3 backend, `infra/`)
- **GitHub Actions** — CI/CD with OIDC auth to AWS
- **Swaggo** — Swagger/OpenAPI 2.0 docs (generated + served)

## Architecture

```
cmd/
├── lambda/          # Lambda entrypoint
│   ├── main.go      # lambda.Start(handler)
│   ├── handler/     # Gin adapter (API Gateway V2)
│   └── provider/    # DI: DB + Gin setup
└── cli/             # Local CLI tool for dev tasks

internal/
├── api/             # HTTP handlers + router setup
├── core/            # Domain models and business logic
└── migrate/         # DB migrations

pkg/
└── dbcon/           # SQLite connection helper

infra/               # Terraform
├── main.tf          # Backend + OIDC provider
├── api-gateway.tf   # HTTP API v2 + routes
├── lambda.tf        # Lambda function + permissions
├── iam.tf           # IAM roles + policies
├── variables.tf     # Reusable variables
├── outputs.tf       # Stack outputs
└── envs/            # tfvars per environment
    ├── dev.tfvars
    └── prod.tfvars
```

## Commands

Run all commands via `make`. See [Makefile](./Makefile).

| Command | Description |
|---|---|
| `make build-cli` | Build the CLI binary |
| `make build-lambda` | Build Lambda bootstrap (linux amd64) |
| `make test` | Run all tests |
| `make swagger` | Regenerate Swagger docs |
| `make tf-plan-dev` | Terraform plan for dev (via Docker) |
| `make tf-plan-prod` | Terraform plan for prod (via Docker) |
| `make tf-apply-dev` | Terraform apply dev (via Docker) |
| `make tf-apply-prod` | Terraform apply prod (via Docker) |
| `make clean` | Clean build artifacts |

## CI/CD (GitHub Actions)

Workflow: `.github/workflows/deploy.yml`

Triggers: push to `main` OR `workflow_dispatch` with environment selection (dev/prod).

Steps:
1. Checkout + OIDC auth to AWS
2. Cache `.terraform/` and Go modules
3. Build Lambda binary
4. `terraform init -reconfigure -backend-config="key=trama/{env}/terraform.tfstate"`
5. `terraform apply -auto-approve`

### Secrets

| Secret | Description |
|---|---|
| `AWS_ACCOUNT_ID` | AWS account ID for OIDC role ARN |

### OIDC IAM Trust Policy

- OIDC provider: `token.actions.githubusercontent.com`
- Audience: `sts.amazonaws.com`
- Action: `sts:AssumeRoleWithWebIdentity`
- Subject: `repo:G-Glez/*` (StringLike)
- Roles: `trama-dev-github`, `trama-prod-github` (per environment)

## Infrastructure

- **Region**: `eu-west-1`
- **State**: S3 `terraform-infra-{account-id}`, key `trama/{env}/terraform.tfstate`
- **API**: HTTP API Gateway v2, routes: `ANY /api/{proxy+}`
- **Lambda**: Go bootstrap binary with Gin adapter
- **CORS**: dev = `*`, prod = `https://trama.app`

## API Docs (OpenAPI / Swagger)

- Generated spec: `docs/swagger.yaml` and `docs/swagger.json`
- Add/update annotations in handler files, then run `make swagger`

## Conventions

- **Language**: Go — standard formatting (`gofumpt`), idiomatic naming
- **Architecture**: domain-driven, hexagonal-style — `internal/core/` for domain, `internal/api/` for delivery
- **Imports**: stdlib first, then third-party, then internal; group with blank lines
- **Errors**: return early, wrap with `fmt.Errorf("context: %w", err)`
- **Testing**: `*_test.go` next to implementation; use `httptest` for HTTP handler tests
- **SQLite**: additive-only migrations in `internal/migrate/migrations/`, executed on Lambda cold start
- **No commented-out code** — delete it
- **No hardcoded secrets** — use env vars loaded via `caarlos0/env`
- **Commits**: [Conventional Commits](https://www.conventionalcommits.org/) — `feat:` for features, `fix:` for bug fixes

## Configuration

Via environment variables (see `.local.env`):

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | Server listen port (local only) |
| `GIN_MODE` | `debug` | Gin mode (`debug`/`release`/`test`) |
| `DATABASE_PATH` | `./data/trama.db` | SQLite database file path |

## Branching Strategy

- **`main`** — production branch. Protected: requires PR, no direct pushes, no force push, linear history.
- **`dev`** — development/integration branch. PRs go here first, then merged to `main`.

Workflow: work on `dev` → PR to `dev` → merge → PR `dev` → `main` → deploy.

`dev` and `main` should stay in sync (merge `main` into `dev` after each deploy).

## Git Policy

**Do not** stage, commit, amend, push, merge, rebase, or create branches/PRs unless the user explicitly grants permission to do so.
