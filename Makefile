# At the top of your Makefile
include .env
export 

run:
	go run ./cmd/api

dev:
	air

build:
	docker build -t gofiber-skeleton .

compose-up:
	docker-compose up

compose-down:
	docker-compose down

sqlc-generate:
	sqlc generate

migrate-up:
	migrate -path db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" -verbose up

migrate-down:
	migrate -path db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" -verbose down