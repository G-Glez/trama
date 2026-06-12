.PHONY: build-api test lint swagger run-api populate-db stop clean

APP_NAME   := trama
BUILD_DIR  := ./build
CMD_DIR    := ./cmd/api

build-api:
	@echo "Building $(APP_NAME)..."
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

test:
	go test ./...

lint:
	golangci-lint run ./...

swagger:
	swag init -g $(CMD_DIR)/main.go -o ./docs

run-api:
	docker compose up --build

-include .env

DATABASE_PATH ?= data/trama.db

populate-db:
	@echo "Clearing and populating database..."
	rm -f "$(DATABASE_PATH)"
	for f in localdb/*.sql; do \
		echo "  Running $$f..."; \
		sqlite3 "$(DATABASE_PATH)" < "$$f"; \
	done
	@echo "Database populated."

stop:
	docker compose down

clean:
	rm -rf $(BUILD_DIR)
