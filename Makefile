.PHONY: help init generate build test test-coverage lint fmt clean docker-build docker-up docker-down migrate-up migrate-down dev run

# Variables
GO := go
GOFLAGS := -v
GOPROXY := direct
BINARY_NAME := service
BINARY_DIR := ./bin

# Help target
help:
	@echo "Available targets:"
	@echo "  init              - Initialize project dependencies"
	@echo "  generate          - Generate sqlc code and mocks"
	@echo "  build             - Build the application"
	@echo "  dev               - Run in development mode with hot reload"
	@echo "  run               - Run the compiled binary"
	@echo "  test              - Run all tests"
	@echo "  test-coverage     - Run tests with coverage report"
	@echo "  lint              - Run golangci-lint"
	@echo "  fmt               - Format code with gofmt"
	@echo "  clean             - Clean build artifacts"
	@echo "  docker-build      - Build Docker image"
	@echo "  docker-up         - Start Docker containers with compose"
	@echo "  docker-down       - Stop Docker containers"
	@echo "  migrate-up        - Run database migrations"
	@echo "  migrate-down      - Rollback database migrations"

# Initialize project
init:
	@echo "Initializing project..."
	$(GO) mod tidy
	$(GO) mod download

# Generate sqlc code and mocks
generate:
	@echo "Generating sqlc code..."
	sqlc generate
	@echo "Generating mocks..."
	$(GO) generate ./...

# Build the application
build: clean
	@echo "Building application..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/service
	@echo "Binary built: $(BINARY_DIR)/$(BINARY_NAME)"

# Run the application
run:
	@echo "Running application..."
	./$(BINARY_DIR)/$(BINARY_NAME)

# Development mode
dev:
	@echo "Running in development mode..."
	$(GO) run ./cmd/service

# Run all tests
test:
	@echo "Running tests..."
	$(GO) test -v -race -count=1 ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter
lint:
	@echo "Running golangci-lint..."
	GOPROXY=$(GOPROXY) golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	gofmt -s -w .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GO) clean
	rm -rf $(BINARY_DIR)
	rm -f coverage.out coverage.html

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t go-fiber-service:latest .

# Start Docker compose
docker-up:
	@echo "Starting Docker containers..."
	docker-compose -f compose.yml up -d

# Stop Docker compose
docker-down:
	@echo "Stopping Docker containers..."
	docker-compose -f compose.yml down

# Run database migrations (requires migrate tool)
migrate-up:
	@echo "Running migrations..."
	migrate -path sql/migrations -database "mysql://root:password@tcp(localhost:3306)/template_go_fiber" up

# Rollback database migrations
migrate-down:
	@echo "Rolling back migrations..."
	migrate -path sql/migrations -database "mysql://root:password@tcp(localhost:3306)/template_go_fiber" down

# Install tools
install-tools:
	@echo "Installing development tools..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	$(GO) install github.com/golang-migrate/migrate/cmd/migrate@latest

.DEFAULT_GOAL := help
