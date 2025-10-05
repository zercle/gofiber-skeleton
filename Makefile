.PHONY: help fmt build run tidy test test-race test-coverage lint sqlc migrate migrate-up migrate-down migrate-create migrate-version generate-docs generate-mocks dev clean ci install-tools

# Default target
help:
	@echo "Available commands:"
	@echo "  make fmt              - Format Go code"
	@echo "  make build            - Build the application"
	@echo "  make run              - Build and run the application"
	@echo "  make dev              - Run with hot-reloading (requires Air)"
	@echo "  make tidy             - Tidy Go modules"
	@echo "  make test             - Run tests"
	@echo "  make test-race        - Run tests with race detector"
	@echo "  make test-coverage    - Run tests with coverage report"
	@echo "  make lint             - Run linter"
	@echo "  make sqlc             - Generate sqlc code"
	@echo "  make migrate-up       - Run database migrations"
	@echo "  make migrate-down     - Rollback last migration"
	@echo "  make migrate-create   - Create new migration (usage: make migrate-create name=create_users)"
	@echo "  make migrate-version  - Show current migration version"
	@echo "  make generate-docs    - Generate Swagger documentation"
	@echo "  make generate-mocks   - Generate mock implementations"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make ci               - Run CI pipeline"
	@echo "  make install-tools    - Install required development tools"

fmt:
	@echo "Formatting Go code..."
	@go fmt ./...

build:
	@echo "Building application..."
	@mkdir -p bin
	@go build -o bin/server cmd/server/main.go

run: build
	@echo "Running application..."
	@./bin/server

dev:
	@echo "Starting development server with hot-reloading..."
	@air

tidy:
	@echo "Tidying Go modules..."
	@go mod tidy

test:
	@echo "Running tests..."
	@go test -v ./...

test-race:
	@echo "Running tests with race detector..."
	@go test -v -race ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint:
	@echo "Running linter..."
	@golangci-lint run --timeout=5m

sqlc:
	@echo "Generating sqlc code..."
	@sqlc generate

migrate-up:
	@echo "Running database migrations..."
	@migrate -path db/migrations -database "$(DATABASE_DSN)" up

migrate-down:
	@echo "Rolling back last migration..."
	@migrate -path db/migrations -database "$(DATABASE_DSN)" down 1

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required. Usage: make migrate-create name=create_users"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	@migrate create -ext sql -dir db/migrations -seq $(name)

migrate-version:
	@echo "Current migration version:"
	@migrate -path db/migrations -database "$(DATABASE_DSN)" version

generate-docs:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go --output ./docs

generate-mocks:
	@echo "Generating mock implementations..."
	@go generate ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf tmp/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

ci: fmt sqlc generate-mocks lint test-race build generate-docs
	@echo "CI pipeline complete"

install-tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Tools installed successfully"
