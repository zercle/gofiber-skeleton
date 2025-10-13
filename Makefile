.PHONY: help build dev test clean migrate-up migrate-down docker-up docker-down docker-logs

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Setup the project (install dependencies and tools)
	@echo "Setting up project..."
	go mod download
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

clean: ## Clean build artifacts
	rm -rf bin/
	go clean -cache

build: clean ## Build the application
	CGO_ENABLED=0 go build -o bin/server cmd/server/main.go

build-prod: clean ## Build the application for production
	CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/server cmd/server/main.go

dev: ## Run the application in development mode with hot-reload
	@echo "Starting development server..."
	air

run: ## Run the application
	go run cmd/server/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	golangci-lint run

lint-fix: ## Run linter with fixes
	golangci-lint run --fix

sqlc: ## Generate Go code from SQL queries
	sqlc generate

schema: ## Generate schema.sql from migration files
	@echo "Generating schema from migrations..."
	./scripts/generate-schema.sh

mocks: ## Generate mocks
	go generate ./...

swagger: ## Generate Swagger documentation
	swag init -g cmd/server/main.go

migrate-up: ## Run database migrations up
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/gofiber_skeleton?sslmode=disable" up

migrate-down: ## Run database migrations down
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/gofiber_skeleton?sslmode=disable" down

migrate-create: ## Create new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: Migration name is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	migrate create -ext sql -dir db/migrations -seq $(NAME)

docker-up: ## Start Docker Compose services
	docker-compose up -d

docker-down: ## Stop Docker Compose services
	docker-compose down

docker-logs: ## Show Docker Compose logs
	docker-compose logs -f

docker-build: ## Build Docker image
	docker build -t gofiber-skeleton .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env gofiber-skeleton

all: setup schema sqlc mocks swagger build test lint ## Run all setup and build tasks