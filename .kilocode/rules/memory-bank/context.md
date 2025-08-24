# Context

**Current Work Focus**
Completed architectural consolidation: centralized SQLC-generated code and SQL queries, integrated swagger annotations for API documentation, configured dependency injection with samber/do, and migrated primary keys to UUIDv7.

**Recent Changes**
- Centralized all SQLC-generated code into `internal/infrastructure/sqlc`.
- Consolidated SQL queries under the `queries` directory.
- Added swagger annotations in each handler `router.go` for API docs.
- Configured samber/do DI container in `internal/infrastructure/app/di.go`.
- Updated configuration loading with Viper in `internal/infrastructure/config/config.go`.
- Applied UUIDv7 functions in database migrations for index-friendly primary keys.

**Next Steps**
1. Regenerate SQLC code (`sqlc generate`) and rebuild the server.
2. Apply database migrations (`migrate -path migrations -database "$DATABASE_URL" up`).
3. Run linting and test suite (`golangci-lint run && go test ./...`).
4. Implement the customer order flow endpoint (`POST /api/v1/orders/create`) and update product stock logic.
5. Update documentation and integration tests for the new endpoint.