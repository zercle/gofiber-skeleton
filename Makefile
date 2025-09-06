.PHONY: build run dev test lint fmt clean migrate migrate-up migrate-down migrate-create docs help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Binary names
MAIN_BINARY=bin/server
MIGRATE_BINARY=bin/migrate

# Build flags
BUILD_FLAGS=-ldflags="-s -w"

## Build the application
build:
	$(GOBUILD) $(BUILD_FLAGS) -o $(MAIN_BINARY) ./cmd/server
	$(GOBUILD) $(BUILD_FLAGS) -o $(MIGRATE_BINARY) ./cmd/migrate

## Run the application
run: build
	./$(MAIN_BINARY)

## Run the application in development mode with hot reloading
dev:
	air

## Run tests
test:
	$(GOTEST) -v -race ./...

## Run tests with coverage
test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## Lint the code
lint:
	golangci-lint run

## Format the code
fmt:
	$(GOFMT) -s -w .
	$(GOCMD) mod tidy

## Generate mocks
generate:
	$(GOCMD) generate ./...

## Run database migrations
migrate: build
	./$(MIGRATE_BINARY)

## Run database migrations up
migrate-up: build
	./$(MIGRATE_BINARY) -direction=up

## Run database migrations down
migrate-down: build
	./$(MIGRATE_BINARY) -direction=down

## Create a new migration file
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

## Generate API documentation
docs:
	swag init -g cmd/server/main.go -o docs

## Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## Install development tools
install-tools:
	$(GOGET) github.com/cosmtrek/air@latest
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/swaggo/swag/cmd/swag@latest
	$(GOGET) github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	$(GOGET) go.uber.org/mock/mockgen@latest

## Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out
	rm -f coverage.html

## Docker commands
docker-build:
	docker build -t gofiber-skeleton .

docker-run:
	docker compose up

docker-stop:
	docker compose down

## Database commands (using docker compose)
db-up:
	docker compose up -d postgres

db-down:
	docker compose down postgres

db-logs:
	docker compose logs -f postgres

## Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Run in development mode with hot reloading"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  lint          - Lint the code"
	@echo "  fmt           - Format the code"
	@echo "  generate      - Generate mocks"
	@echo "  migrate       - Run database migrations"
	@echo "  migrate-up    - Run database migrations up"
	@echo "  migrate-down  - Run database migrations down"
	@echo "  migrate-create- Create a new migration file"
	@echo "  docs          - Generate API documentation"
	@echo "  deps          - Install dependencies"
	@echo "  install-tools - Install development tools"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-stop   - Stop Docker Compose"
	@echo "  db-up         - Start database with Docker"
	@echo "  db-down       - Stop database"
	@echo "  db-logs       - Show database logs"
	@echo "  help          - Show this help message"