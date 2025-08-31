include .env
export

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@egrep '^(.+)\s*:.*##\s*(.+)' $(MAKEFILE_LIST) | column -t -c 2 -s ':#'

.PHONY: setup
setup: ## Set up the project
	@echo "Setting up the project..."
	@go mod download
	@go mod tidy

.PHONY: dev
dev: ## Run the application in development mode with hot reload
	@echo "Starting development server with hot reload..."
	@air

.PHONY: build
build: ## Build the application
	@echo "Building application..."
	@go build -o bin/server cmd/server/main.go

.PHONY: run
run: ## Run the application
	@echo "Running application..."
	@./bin/server

.PHONY: generate
generate: ## Generate mocks and other code
	@echo "Generating code..."
	@go generate ./...

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	@go clean -testcache
	@go test -v -race ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run --fix ./...

.PHONY: lint-fix
lint-fix: ## Run linter with auto-fix
	@echo "Running linter with auto-fix..."
	@golangci-lint run --fix ./...

.PHONY: check
check: generate lint test ## Run all checks (generate, lint, test)
	@echo "All checks completed successfully!"

.PHONY: swagger-generate
swagger-generate: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go -o docs

.PHONY: sqlc-generate
sqlc-generate: ## Generate SQLC code
	@echo "Generating SQLC code..."
	@sqlc generate

.PHONY: migrate-up
migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	@migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up

.PHONY: migrate-down
migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	@migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down

.PHONY: migrate-create
migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@echo "Creating migration: $(name)"
	@migrate create -ext sql -dir db/migrations -seq $(name)

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t gofiber-skeleton .

.PHONY: docker-up
docker-up: ## Start services with Docker Compose
	@echo "Starting services with Docker Compose..."
	@docker compose up -d

.PHONY: docker-down
docker-down: ## Stop services with Docker Compose
	@echo "Stopping services with Docker Compose..."
	@docker compose down

.PHONY: docker-logs
docker-logs: ## Show Docker Compose logs
	@docker compose logs -f

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/ tmp/
	@rm -f coverage.out coverage.html

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install go.uber.org/mock/mockgen@latest