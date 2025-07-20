DB_URL=postgres://user:password@db:5432/gofiber?sslmode=disable

.PHONY: dev
dev:
	air

.PHONY: migrate-up
migrate-up:
	docker-compose exec app migrate -path db/migrations -database "$(DB_URL)" -verbose up

.PHONY: migrate-down
migrate-down:
	docker-compose exec app migrate -path db/migrations -database "$(DB_URL)" -verbose down

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: test
test:
	go test ./...

.PHONY: swag
swag:
	swag init -g cmd/api/main.go

.PHONY: mockgen
mockgen:
	mockgen -source=internal/usecases/interfaces.go -destination=mocks/usecases.go -package=mocks
