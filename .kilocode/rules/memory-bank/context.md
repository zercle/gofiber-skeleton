# Context

**Current Work Focus**
Enabled SQLC `emit_methods_with_db_argument: true` for repository-level transaction control, replaced generated `null` types with `github.com/guregu/null/v6`, and implemented demo endpoints for full transaction flow and complex SQL joins.

**Recent Changes**
- Centralized all SQLC-generated code into `internal/infrastructure/sqlc`.
- Consolidated SQL queries under the `queries` directory.
- Added swagger annotations in each handler `router.go` for API docs.
- Configured samber/do DI container in `internal/infrastructure/app/di.go`.
- Updated configuration loading with Viper in `internal/infrastructure/config/config.go`.
- Applied UUIDv7 functions in database migrations for index-friendly primary keys.

**Next Steps**
1.  Update SQLC configuration (`emit_methods_with_db_argument: true`) and regenerate code (`sqlc generate`).
2.  Apply database migrations (`migrate -path migrations -database "$DATABASE_URL" up`).
3.  Run linting and test suite (`golangci-lint run && go test ./...`).
4.  Implement demo transaction endpoint (`POST /api/v1/demo/transaction`) and ensure all repository operations use transactions.
5.  Implement demo join endpoint (`GET /api/v1/demo/joined`) with advanced SQL join query and write integration tests.
6.  Update documentation and memory bank files to reflect changes.