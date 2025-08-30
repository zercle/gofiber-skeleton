# Context

**Current Work Focus**

Refactoring domain interfaces into each domain package, relocating database migrations and SQL query files under `db/migrations` and `db/queries`, and ensuring all database interactions are testable without a real database.

**Recent Changes**

- Removed demo transaction and demo join endpoints.
- Updated memory bank instructions to focus on core CRUD stories and multi-stage query guidance.
- Consolidated multi-stage join query examples into product, order, and user repository documentation.
- Adjusted SQLC configuration for `emit_methods_with_db_argument` to support multi-stage queries in domain repositories.

**Next Steps**

1. Regenerate SQLC code (`sqlc generate`) after updating query paths for code generation. (Completed)
2. Relocate `migrations` and `queries` directories under `db/` and adjust any path references. (Completed)
3. Update tests to use `go-sqlmock` for all database interactions (no real DB connections).
4. Run linting and unit tests (`golangci-lint run && go test ./...`).
5. Update documentation and memory bank files to reflect interface relocations and folder restructures.