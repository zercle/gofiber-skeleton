.PHONY: help build run test clean generate-mocks docker-build docker-run migrate-up migrate-down sqlc-generate lint fmt deps install-tools

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application locally"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  generate-mocks - Generate mock files"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  sqlc-generate - Generate SQLC code"
	@echo "  lint          - Lint code"
	@echo "  fmt           - Format code"
	@echo "  deps          - Install dependencies"
	@echo "  install-tools - Install development tools"
	@echo "  swag          - Generate API document"

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application locally
run:
	go run ./cmd/server

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Generate mock files
generate-mocks:
	go generate ./...

# Build Docker image
docker-build:
	docker build -t ecommerce-api .

# Run with Docker Compose
docker-run:
	docker compose up --build

# Run database migrations
migrate-up:
	migrate -path db/migrations -database "postgres://postgres:password@localhost:5432/ecommerce?sslmode=disable" up

# Rollback database migrations
migrate-down:
	migrate -path db/migrations -database "postgres://postgres:password@localhost:5432/ecommerce?sslmode=disable" down

# Install dependencies
deps:
	go mod download
	go mod tidy

# Generate SQLC code
sqlc-generate:
	sqlc generate

# Lint code
lint:
	golangci-lint run --fix

# Format code
fmt:
	go fmt ./...
	goimports -w .

swag:
	swag init --generalInfo ./cmd/server/main.go --parseDependency --parseInternal --parseFuncBody

# Install development tools
install-tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/cosmtrek/air@latest