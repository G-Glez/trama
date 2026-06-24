.PHONY: build-cli build-lambda test lint swagger clean               \
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

TF_IMAGE    := hashicorp/terraform:1.15.6
TF_TMPDIR   := /tmp/opencode/trama-infra
TF_CACHEDIR := /tmp/opencode/trama-tf-cache
TF_ENV      := -e AWS_PROFILE -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -e AWS_SESSION_TOKEN
TF_DOCKER   := docker run --rm --entrypoint sh \
                -v $(TF_TMPDIR)/infra:/workspace/infra \
                -v $(TF_TMPDIR)/build:/workspace/build \
                -v $(HOME)/.aws:/root/.aws:ro \
                -v $(TF_CACHEDIR):/workspace/infra/.terraform \
                $(TF_ENV) -w /workspace/infra $(TF_IMAGE)

define TF_INIT
terraform init -reconfigure -backend-config="key=trama/$(1)/terraform.tfstate"
endef

define TF_IMPORT_DEV
terraform import -var-file=envs/dev.tfvars aws_dynamodb_table.users trama-users-dev || true && \
terraform import -var-file=envs/dev.tfvars aws_iam_role.lambda trama-lambda-dev || true && \
terraform import -var-file=envs/dev.tfvars aws_iam_role.github trama-github-dev || true && \
terraform import -var-file=envs/dev.tfvars aws_iam_policy.logs arn:aws:iam::211125667058:policy/trama-logs-dev || true && \
terraform import -var-file=envs/dev.tfvars aws_iam_policy.dynamodb arn:aws:iam::211125667058:policy/trama-dynamodb-dev || true && \
terraform import -var-file=envs/dev.tfvars aws_cloudwatch_log_group.trama /aws/lambda/trama-api-dev || true
endef

define TF_PREPARE
rm -rf $(TF_TMPDIR) && \
mkdir -p $(TF_TMPDIR)/infra $(TF_TMPDIR)/build/lambda && \
cp -r infra/* $(TF_TMPDIR)/infra/ && \
cp build/lambda/bootstrap $(TF_TMPDIR)/build/lambda/
endef

tf-plan-dev: build-lambda       ## Plan dev (via Docker)
	$(TF_PREPARE) && \
	$(TF_DOCKER) -c '$(call TF_INIT,dev) && terraform plan -var-file=envs/dev.tfvars' && \
	rm -rf $(TF_TMPDIR)

tf-plan-prod: build-lambda      ## Plan prod (via Docker)
	$(TF_PREPARE) && \
	$(TF_DOCKER) -c '$(call TF_INIT,prod) && terraform plan -var-file=envs/prod.tfvars' && \
	rm -rf $(TF_TMPDIR)

tf-apply-dev: build-lambda       ## Apply dev (via Docker)
	$(TF_PREPARE) && \
	$(TF_DOCKER) -c '$(call TF_INIT,dev) && ($(TF_IMPORT_DEV)) && terraform apply -auto-approve -var-file=envs/dev.tfvars' && \
	rm -rf $(TF_TMPDIR)

tf-apply-prod: build-lambda      ## Apply prod (via Docker)
	$(TF_PREPARE) && \
	$(TF_DOCKER) -c '$(call TF_INIT,prod) && terraform apply -auto-approve -var-file=envs/prod.tfvars' && \
	rm -rf $(TF_TMPDIR)

tf-destroy-dev: build-lambda     ## Destroy dev (via Docker)
	$(TF_PREPARE) && \
	$(TF_DOCKER) -c '$(call TF_INIT,dev) && terraform destroy -auto-approve -var-file=envs/dev.tfvars' && \
	rm -rf $(TF_TMPDIR)

tf-destroy-prod: build-lambda    ## Destroy prod (via Docker)
	$(TF_PREPARE) && \
	$(TF_DOCKER) -c '$(call TF_INIT,prod) && terraform destroy -auto-approve -var-file=envs/prod.tfvars' && \
	rm -rf $(TF_TMPDIR)

## ──────────────────────────────
## Cleanup
## ──────────────────────────────

clean:                          ## Remove build artifacts
	rm -rf $(BUILD_DIR)
