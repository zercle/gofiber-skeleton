.PHONY: help build run test test-coverage clean dev migrate-up migrate-down migrate-create sqlc lint format tidy docker-build docker-up docker-down generate-mocks docs

# Variables
APP_NAME=gofiber-skeleton
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

## help: Display this help message
help:
	@echo "Available commands:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the application binary
build:
	@echo "Building ${APP_NAME}..."
	@go build ${LDFLAGS} -o bin/${APP_NAME} cmd/server/main.go

## run: Run the application
run:
	@echo "Running ${APP_NAME}..."
	@go run cmd/server/main.go

## dev: Run with hot reload using Air
dev:
	@echo "Starting development server with hot reload..."
	@air

## test: Run all tests
test:
	@echo "Running tests..."
	@go test -v -race ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/ tmp/ coverage.out coverage.html

## migrate-up: Run database migrations up
migrate-up:
	@echo "Running migrations..."
	@migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/gofiber_skeleton?sslmode=disable" up

## migrate-down: Rollback last migration
migrate-down:
	@echo "Rolling back migration..."
	@migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/gofiber_skeleton?sslmode=disable" down 1

## migrate-create: Create a new migration file (usage: make migrate-create NAME=create_users)
migrate-create:
	@echo "Creating migration: ${NAME}..."
	@migrate create -ext sql -dir db/migrations -seq ${NAME}

## sqlc: Generate Go code from SQL
sqlc:
	@echo "Generating sqlc code..."
	@sqlc generate

## generate-mocks: Generate mock implementations
generate-mocks:
	@echo "Generating mocks..."
	@go generate ./...

## docs: Generate Swagger documentation
docs:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/server/main.go -o docs

## lint: Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

## format: Format code
format:
	@echo "Formatting code..."
	@gofmt -s -w .
	@goimports -w .

## tidy: Tidy and verify dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@go mod verify

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t ${APP_NAME}:${VERSION} -t ${APP_NAME}:latest .

## docker-up: Start Docker Compose services
docker-up:
	@echo "Starting Docker services..."
	@docker-compose up -d

## docker-down: Stop Docker Compose services
docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down

## setup: Initial project setup
setup:
	@echo "Setting up project..."
	@go mod download
	@go install github.com/cosmtrek/air@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install go.uber.org/mock/mockgen@latest
	@cp .env.example .env
	@echo "Setup complete! Edit .env file with your configuration."

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install go.uber.org/mock/mockgen@latest
	@echo "Tools installed successfully!"
