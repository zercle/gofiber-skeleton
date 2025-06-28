APP_NAME=gofiber-boilerplate
BUILD_DIR=./bin

.PHONY: all build run test clean migrate-up migrate-down docker-build docker-up

all: build

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/app
	@echo "Build complete."

run:
	@echo "Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

test:
	@echo "Running tests..."
	@PWD=$(CURDIR) go test ./...

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -f *.db
	@echo "Clean complete."

migrate-up:
	@echo "Running database migrations up..."
	@migrate -path database/migrations -database "sqlite3://$(shell grep DATABASE_URL config/GO_ENV.yaml | awk '{print $2}')" up
	@echo "Migrations up complete."

migrate-down:
	@echo "Running database migrations down..."
	@migrate -path database/migrations -database "sqlite3://$(shell grep DATABASE_URL config/GO_ENV.yaml | awk '{print $2}')" down
	@echo "Migrations down complete."

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME) .
	@echo "Docker image built."

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Docker containers started."


generate-proto:
	@echo "Generating protobuf code..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/user/user.proto api/product/product.proto api/order/order.proto
	@echo "Protobuf code generated."

generate-mocks:
	@echo "Generating mocks..."
	@go generate ./...
	@echo "Mocks generated."


