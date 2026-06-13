.PHONY: build-api build-cli test lint swagger sqlc run-api populate-db migrate set-env stop clean

APP_NAME      := trama
CLI_NAME      := trama-cli
BUILD_DIR     := ./build
CMD_DIR       := ./cmd/api
CLI_DIR       := ./cmd/cli
-include .local.env

.EXPORT_ALL_VARIABLES:

build-api:
	@echo "Building $(APP_NAME)..."
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

build-cli:
	@echo "Building $(CLI_NAME)..."
	go build -o $(BUILD_DIR)/$(CLI_NAME) $(CLI_DIR)

test:
	go test ./...

lint:
	golangci-lint run ./...

swagger:
	swag init -g ./internal/api/router.go -o ./docs

sqlc:
	sqlc generate

run-api:
	docker compose up --build

populate-db:
	@echo "Clearing and populating database..."
	rm -f "$(DATABASE_PATH)"
	for f in localdb/*.sql; do \
		echo "  Running $$f..."; \
		sqlite3 "$(DATABASE_PATH)" < "$$f"; \
	done
	@echo "Database populated."

migrate: build-cli
	./build/$(CLI_NAME) migrate

set-env:
	@echo "=== Environment loaded from .local.env ==="
	@echo "  DATABASE_PATH=$(DATABASE_PATH)"
	@echo ""
	@echo "Variables are exported to subprocesses. Run:"
	@echo "  make migrate"

stop:
	docker compose down

clean:
	rm -rf $(BUILD_DIR)
