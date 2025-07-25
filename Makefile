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
	migrate -path db/migrations -database "postgres://user:password@localhost:5432/shortener?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migrations -database "postgres://user:password@localhost:5432/shortener?sslmode=disable" -verbose down
