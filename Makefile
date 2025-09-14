.PHONY: fmt build run tidy migrate test

fmt:
	go fmt ./...

build:
	go build -o bin/server cmd/server/main.go

run: build
	./bin/server

tidy:
	go mod tidy

migrate:
	@echo "Placeholder for migration tool execution"
	@echo "Example: migrate -path migrations -database "postgres://user:password@host:port/database?sslmode=disable" up"

test:
	go test ./...

generate-docs:
	swag init -g cmd/server/main.go --output ./docs

sqlc:
	sqlc generate

lint:
	golangci-lint run

test-race:
	go test ./... -race

ci: fmt sqlc lint test-race build generate-docs