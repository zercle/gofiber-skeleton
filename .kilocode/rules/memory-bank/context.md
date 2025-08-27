# Context

**Current Work Focus**

Simplifying boilerplate by removing standalone demo endpoints and integrating guidance for complex multi-stage queries within domain repositories.

**Recent Changes**

- Removed demo transaction and demo join endpoints.
- Updated memory bank instructions to focus on core CRUD stories and multi-stage query guidance.
- Consolidated multi-stage join query examples into product, order, and user repository documentation.
- Adjusted SQLC configuration for `emit_methods_with_db_argument` to support multi-stage queries in domain repositories.

**Next Steps**

1. Regenerate SQLC code (`sqlc generate`) after adjusting configuration.
2. Apply database migrations (`migrate -path migrations -database "$DATABASE_URL" up`).
3. Run linting and tests (`golangci-lint run && go test ./...`).
4. Implement multi-stage join queries in repository implementations and add example code.
5. Write integration tests for multi-stage query functionality.
6. Update documentation and memory bank files to reflect changes.