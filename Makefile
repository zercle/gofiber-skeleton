.PHONY: help dev build test lint clean migrate-up migrate-down migrate-create sqlc swag mocks docker-up docker-down

# Variables
APP_NAME := gofiber-skeleton
MAIN_PATH := ./cmd/server
MIGRATE_PATH := ./cmd/migrate
MIGRATION_DIR := ./db/migrations
BUILD_DIR := ./bin
DOCKER_COMPOSE := docker-compose

# Colors for output
COLOR_RESET := \033[0m
COLOR_BLUE := \033[34m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m

##@ General

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

dev: ## Start development server with hot reload
	@printf "$(COLOR_BLUE)Starting development server...$(COLOR_RESET)\n"
	@air -c .air.toml

run: ## Run the application
	@printf "$(COLOR_BLUE)Running application...$(COLOR_RESET)\n"
	@go run $(MAIN_PATH)/main.go

##@ Build

build: ## Build the application
	@printf "$(COLOR_BLUE)Building $(APP_NAME)...$(COLOR_RESET)\n"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@printf "$(COLOR_GREEN)Build complete: $(BUILD_DIR)/$(APP_NAME)$(COLOR_RESET)\n"

build-linux: ## Build for Linux
	@printf "$(COLOR_BLUE)Building for Linux...$(COLOR_RESET)\n"
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux $(MAIN_PATH)
	@printf "$(COLOR_GREEN)Build complete: $(BUILD_DIR)/$(APP_NAME)-linux$(COLOR_RESET)\n"

clean: ## Clean build artifacts
	@printf "$(COLOR_YELLOW)Cleaning build artifacts...$(COLOR_RESET)\n"
	@rm -rf $(BUILD_DIR)
	@rm -rf tmp
	@printf "$(COLOR_GREEN)Clean complete$(COLOR_RESET)\n"

##@ Testing

test: ## Run tests
	@printf "$(COLOR_BLUE)Running tests...$(COLOR_RESET)\n"
	@go test -v -race -coverprofile=coverage.out ./...
	@printf "$(COLOR_GREEN)Tests complete$(COLOR_RESET)\n"

test-coverage: ## Run tests with coverage report
	@printf "$(COLOR_BLUE)Running tests with coverage...$(COLOR_RESET)\n"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@printf "$(COLOR_GREEN)Coverage report: coverage.html$(COLOR_RESET)\n"

test-unit: ## Run unit tests only
	@printf "$(COLOR_BLUE)Running unit tests...$(COLOR_RESET)\n"
	@go test -v -short ./...

##@ Code Quality

lint: ## Run linter
	@printf "$(COLOR_BLUE)Running linter...$(COLOR_RESET)\n"
	@golangci-lint run ./...
	@printf "$(COLOR_GREEN)Linting complete$(COLOR_RESET)\n"

fmt: ## Format code
	@printf "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)\n"
	@go fmt ./...
	@printf "$(COLOR_GREEN)Formatting complete$(COLOR_RESET)\n"

vet: ## Run go vet
	@printf "$(COLOR_BLUE)Running go vet...$(COLOR_RESET)\n"
	@go vet ./...
	@printf "$(COLOR_GREEN)Vet complete$(COLOR_RESET)\n"

##@ Database

migrate-up: ## Run database migrations
	@printf "$(COLOR_BLUE)Running migrations...$(COLOR_RESET)\n"
	@go run $(MIGRATE_PATH)/main.go up
	@printf "$(COLOR_GREEN)Migrations complete$(COLOR_RESET)\n"

migrate-down: ## Rollback last migration
	@printf "$(COLOR_YELLOW)Rolling back migration...$(COLOR_RESET)\n"
	@go run $(MIGRATE_PATH)/main.go down
	@printf "$(COLOR_GREEN)Rollback complete$(COLOR_RESET)\n"

migrate-force: ## Force migration version (use: make migrate-force VERSION=1)
	@printf "$(COLOR_YELLOW)Forcing migration to version $(VERSION)...$(COLOR_RESET)\n"
	@go run $(MIGRATE_PATH)/main.go force $(VERSION)
	@printf "$(COLOR_GREEN)Force complete$(COLOR_RESET)\n"

migrate-version: ## Show current migration version
	@printf "$(COLOR_BLUE)Checking migration version...$(COLOR_RESET)\n"
	@go run $(MIGRATE_PATH)/main.go version

migrate-create: ## Create new migration (use: make migrate-create NAME=create_users_table)
	@printf "$(COLOR_BLUE)Creating migration: $(NAME)...$(COLOR_RESET)\n"
	@migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(NAME)
	@printf "$(COLOR_GREEN)Migration created$(COLOR_RESET)\n"

##@ Code Generation

sqlc: ## Generate SQL code
	@printf "$(COLOR_BLUE)Generating SQL code...$(COLOR_RESET)\n"
	@sqlc generate
	@printf "$(COLOR_GREEN)SQL code generated$(COLOR_RESET)\n"

swag: ## Generate API documentation
	@printf "$(COLOR_BLUE)Generating API documentation...$(COLOR_RESET)\n"
	@swag init -g $(MAIN_PATH)/main.go -o ./docs
	@printf "$(COLOR_GREEN)API documentation generated$(COLOR_RESET)\n"

mocks: ## Generate mocks for testing
	@printf "$(COLOR_BLUE)Generating mocks...$(COLOR_RESET)\n"
	@go generate ./...
	@printf "$(COLOR_GREEN)Mocks generated$(COLOR_RESET)\n"

##@ Docker

docker-up: ## Start Docker containers
	@printf "$(COLOR_BLUE)Starting Docker containers...$(COLOR_RESET)\n"
	@$(DOCKER_COMPOSE) up -d
	@printf "$(COLOR_GREEN)Containers started$(COLOR_RESET)\n"

docker-down: ## Stop Docker containers
	@printf "$(COLOR_YELLOW)Stopping Docker containers...$(COLOR_RESET)\n"
	@$(DOCKER_COMPOSE) down
	@printf "$(COLOR_GREEN)Containers stopped$(COLOR_RESET)\n"

docker-logs: ## Show Docker logs
	@$(DOCKER_COMPOSE) logs -f

docker-build: ## Build Docker image
	@printf "$(COLOR_BLUE)Building Docker image...$(COLOR_RESET)\n"
	@docker build -t $(APP_NAME):latest .
	@printf "$(COLOR_GREEN)Docker image built$(COLOR_RESET)\n"

##@ Dependencies

deps: ## Download dependencies
	@printf "$(COLOR_BLUE)Downloading dependencies...$(COLOR_RESET)\n"
	@go mod download
	@printf "$(COLOR_GREEN)Dependencies downloaded$(COLOR_RESET)\n"

deps-update: ## Update dependencies
	@printf "$(COLOR_BLUE)Updating dependencies...$(COLOR_RESET)\n"
	@go get -u ./...
	@go mod tidy
	@printf "$(COLOR_GREEN)Dependencies updated$(COLOR_RESET)\n"

deps-verify: ## Verify dependencies
	@printf "$(COLOR_BLUE)Verifying dependencies...$(COLOR_RESET)\n"
	@go mod verify
	@printf "$(COLOR_GREEN)Dependencies verified$(COLOR_RESET)\n"

##@ Tools

install-tools: ## Install development tools
	@printf "$(COLOR_BLUE)Installing development tools...$(COLOR_RESET)\n"
	@go install github.com/cosmtrek/air@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install go.uber.org/mock/mockgen@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@printf "$(COLOR_GREEN)Tools installed$(COLOR_RESET)\n"
