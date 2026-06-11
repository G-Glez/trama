.PHONY: build-api test lint swagger run-api stop clean

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

stop:
	docker compose down

clean:
	rm -rf $(BUILD_DIR)
