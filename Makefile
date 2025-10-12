.PHONY: help build dev test lint clean migrate-up migrate-down sqlc swag mocks docker-up docker-down

# Variables
APP_NAME := gofiber-skeleton
BUILD_DIR := build
DOCKER_COMPOSE := docker-compose

# Go variables
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development targets
dev: ## Start development server with hot reload
	@if command -v air >/dev/null 2>&1; then \
		air -c .air.toml; \
	else \
		echo "Air not found. Installing air..."; \
		$(GOGET) -u github.com/cosmtrek/air@latest; \
		air -c .air.toml; \
	fi

run: ## Run the application without hot reload
	$(GOCMD) run cmd/server/main.go

# Build targets
build: ## Build the application
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) cmd/server/main.go

build-all: ## Build the application for multiple platforms
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 cmd/server/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 cmd/server/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe cmd/server/main.go

# Test targets
test: ## Run all tests
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

test-unit: ## Run unit tests only
	$(GOTEST) -v ./internal/... ./pkg/...

test-integration: ## Run integration tests
	$(GOTEST) -v -tags=integration ./...

# Quality targets
lint: ## Run linter
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Installing..."; \
		$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

fmt: ## Format Go code
	$(GOFMT) -s -w .

fmt-check: ## Check if code is formatted
	$(GOFMT) -s -d .

vet: ## Run go vet
	$(GOCMD) vet ./...

tidy: ## Tidy go modules
	$(GOMOD) tidy

mod-verify: ## Verify go modules
	$(GOMOD) verify

# Database targets
migrate-up: ## Run database migrations up
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up; \
	else \
		echo "migrate tool not found. Installing..."; \
		$(GOGET) -u github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
		migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up; \
	fi

migrate-down: ## Run database migrations down
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down; \
	else \
		echo "migrate tool not found. Installing..."; \
		$(GOGET) -u github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
		migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down; \
	fi

migrate-create: ## Create new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME parameter is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@if command -v migrate >/dev/null 2>&1; then \
		migrate create -ext sql -dir db/migrations -seq $(NAME); \
	else \
		echo "migrate tool not found. Installing..."; \
		$(GOGET) -u github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
		migrate create -ext sql -dir db/migrations -seq $(NAME); \
	fi

# SQLC targets
sqlc: ## Generate SQLC code
	@if command -v sqlc >/dev/null 2>&1; then \
		sqlc generate; \
	else \
		echo "sqlc not found. Installing..."; \
		$(GOGET) -u github.com/sqlc-dev/sqlc/cmd/sqlc@latest; \
		sqlc generate; \
	fi

sqlc-verify: ## Verify SQLC generated code is up to date
	@if command -v sqlc >/dev/null 2>&1; then \
		sqlc diff; \
	else \
		echo "sqlc not found. Installing..."; \
		$(GOGET) -u github.com/sqlc-dev/sqlc/cmd/sqlc@latest; \
		sqlc diff; \
	fi

# Swagger targets
swag: ## Generate Swagger documentation
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g cmd/server/main.go -o docs; \
	else \
		echo "swag not found. Installing..."; \
		$(GOGET) -u github.com/swaggo/swag/cmd/swag@latest; \
		swag init -g cmd/server/main.go -o docs; \
	fi

# Mock targets
mocks: ## Generate mocks
	@if command -v mockgen >/dev/null 2>&1; then \
		find . -name "*.go" -type f -exec grep -l "//go:generate" {} \; | xargs -n 1 dirname | sort -u | while read dir; do \
			cd "$$dir" && go generate ./...; cd -; \
		done; \
	else \
		echo "mockgen not found. Installing..."; \
		$(GOGET) -u github.com/golang/mock/mockgen@latest; \
		find . -name "*.go" -type f -exec grep -l "//go:generate" {} \; | xargs -n 1 dirname | sort -u | while read dir; do \
			cd "$$dir" && go generate ./...; cd -; \
		done; \
	fi

# Docker targets
docker-up: ## Start Docker containers
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop Docker containers
	$(DOCKER_COMPOSE) down

docker-logs: ## Show Docker logs
	$(DOCKER_COMPOSE) logs -f

docker-build: ## Build Docker image
	docker build -t $(APP_NAME) .

# Setup targets
setup: ## Setup development environment
	@echo "Setting up development environment..."
	$(GOMOD) download
	$(GOGET) -u github.com/cosmtrek/air@latest
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) -u github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	$(GOGET) -u github.com/swaggo/swag/cmd/swag@latest
	$(GOGET) -u github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	$(GOGET) -u github.com/golang/mock/mockgen@latest
	cp .env.example .env
	@echo "Development environment setup complete!"
	@echo "Please update the .env file with your configuration."

setup-ci: ## Setup CI environment
	@echo "Setting up CI environment..."
	$(GOMOD) download
	sqlc generate
	swag init -g cmd/server/main.go -o docs
	@echo "CI environment setup complete!"

# Clean targets
clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

clean-docker: ## Clean Docker resources
	docker system prune -f
	docker volume prune -f

# Git hooks
install-hooks: ## Install git hooks
	@echo "Installing git hooks..."
	cp scripts/pre-commit .git/hooks/
	chmod +x .git/hooks/pre-commit
	cp scripts/pre-push .git/hooks/
	chmod +x .git/hooks/pre-push
	@echo "Git hooks installed!"

# Quick start
quick-start: ## Quick start development (setup + docker-up + migrate-up)
	make setup
	make docker-up
	sleep 5
	make migrate-up
	make sqlc
	make swag
	make mocks
	@echo "Quick start complete! Run 'make dev' to start the development server."

# CI pipeline
ci: lint test-coverage ## Run CI pipeline

# Production build
release: clean fmt lint test-coverage build-all ## Build for production