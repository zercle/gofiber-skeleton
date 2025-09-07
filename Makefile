.PHONY: help build run dev test clean migrate-up migrate-down migrate-create docker-up docker-down sqlc swagger

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Run with hot reload using Air"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  migrate-up    - Run database migrations up"
	@echo "  migrate-down  - Run database migrations down"
	@echo "  migrate-create - Create new migration file"
	@echo "  docker-up     - Start Docker services"
	@echo "  docker-down   - Stop Docker services"
	@echo "  sqlc          - Generate sqlc code"
	@echo "  swagger       - Generate Swagger documentation"
	@echo "  lint          - Run golangci-lint"
	@echo "  fmt           - Format code"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/server/main.go

# Development with hot reload
dev:
	@echo "Starting development server with hot reload..."
	air

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out coverage.html
	rm -f air.log build-errors.log

# Database migrations
migrate-up:
	@echo "Running database migrations up..."
	go run cmd/migrate/main.go -direction=up

migrate-down:
	@echo "Running database migrations down..."
	go run cmd/migrate/main.go -direction=down -steps=1

migrate-create:
	@echo "Creating new migration..."
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Docker commands
docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose logs -f

# Code generation
sqlc:
	@echo "Generating sqlc code..."
	sqlc generate

swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/server/main.go -o docs

# Code quality
lint:
	@echo "Running golangci-lint..."
	golangci-lint run

fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Setup development environment
setup:
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file from .env.example"; fi
	@echo "Please update .env file with your configuration"
	@echo "Run 'make docker-up' to start required services"

# All-in-one development setup
dev-setup: setup docker-up
	@echo "Waiting for services to be ready..."
	@sleep 5
	@make migrate-up
	@echo "Development environment ready!"
	@echo "Run 'make dev' to start the development server"

# Production build
build-prod:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/server/main.go

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t gofiber-skeleton:latest .