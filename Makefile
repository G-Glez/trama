.PHONY: build-cli build-lambda test lint swagger clean               \
        tf-init tf-docker tf-docker-cmd                              \
        tf-plan-dev tf-plan-prod tf-apply-dev tf-apply-prod

## Project
CLI_NAME      := trama-cli
BUILD_DIR     := ./build
CLI_DIR       := ./cmd/cli
-include .local.env

.EXPORT_ALL_VARIABLES:

default: help

help:                           ## Show this help
	@grep -Eh '^[a-z].+:.*##' $(MAKEFILE_LIST) | sort | \
	  awk 'BEGIN {FS = ":.*##"}; {printf "  \033[1;34m%-20s\033[0m %s\n", $$1, $$2}'

## ──────────────────────────────
## Go builds
## ──────────────────────────────

build-cli:                      ## Build the CLI binary
	@echo "Building $(CLI_NAME)..."
	go build -o $(BUILD_DIR)/$(CLI_NAME) $(CLI_DIR)

build-lambda:                   ## Build Lambda bootstrap binary (linux amd64)
	@echo "Building Lambda bootstrap..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/lambda/bootstrap ./cmd/lambda/
	@echo "Lambda binary built at $(BUILD_DIR)/lambda/bootstrap"

## ──────────────────────────────
## Quality
## ──────────────────────────────

test:                           ## Run all tests
	go test ./...

lint:                           ## Run linter
	golangci-lint run ./...

## ──────────────────────────────
## Code generation
## ──────────────────────────────

swagger:                        ## Regenerate Swagger/OpenAPI docs
	@which swag >/dev/null 2>&1 || (echo "Installing swag..."; go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g ./internal/api/router.go -o ./docs

## ──────────────────────────────
## AWS / Terraform
## ──────────────────────────────

TF_IMAGE   := hashicorp/terraform:1.15.6
TF_VOLUMES := -v $(PWD):/workspace -v $(HOME)/.aws:/root/.aws:ro
TF_ENV     := -e AWS_PROFILE=$(AWS_PROFILE)
TF_RUN     := docker run --rm -it $(TF_VOLUMES) $(TF_ENV)

tf-init-dev:                    ## Init S3 backend for dev
	$(TF_RUN) -w /workspace/infra $(TF_IMAGE) init -backend-config="key=trama/dev/terraform.tfstate"

tf-init-prod:                   ## Init S3 backend for prod
	$(TF_RUN) -w /workspace/infra $(TF_IMAGE) init -backend-config="key=trama/prod/terraform.tfstate"

tf-plan-dev: build-lambda       ## Plan dev (via Docker)
	$(TF_RUN) -w /workspace/infra $(TF_IMAGE) plan -var-file=envs/dev.tfvars

tf-plan-prod: build-lambda      ## Plan prod (via Docker)
	$(TF_RUN) -w /workspace/infra $(TF_IMAGE) plan -var-file=envs/prod.tfvars

tf-apply-dev: build-lambda      ## Apply dev (via Docker)
	$(TF_RUN) -w /workspace/infra $(TF_IMAGE) apply -var-file=envs/dev.tfvars

tf-apply-prod: build-lambda     ## Apply prod (via Docker)
	$(TF_RUN) -w /workspace/infra $(TF_IMAGE) apply -var-file=envs/prod.tfvars

## ──────────────────────────────
## Cleanup
## ──────────────────────────────

clean:                          ## Remove build artifacts
	rm -rf $(BUILD_DIR)
