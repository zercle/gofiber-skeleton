# Context

**Current Work Focus**
Finalizing the architectural update to centralize SQLC-generated code and query files, integrate API documentation via gofiber/swagger, and configure application dependency injection using samber/do.

**Recent Changes**
- Centralized all SQLC-generated code into the `internal/infrastructure/sqlc` directory.
- Moved all SQL query files to a single `queries` directory at the project root.
- Updated `sqlc.yaml` to reflect the new centralized structure.
- Updated `architecture.md` and `tech.md` to document the new SQLC configuration and file locations.
- Defined interfaces for each domain layer in `internal/domain`.
- Added swagger annotations requirement to each handler `router.go` to generate API documentation.

**Next Steps**
1. Change entry point for each domain from `router.go` to `router.go`, implement dependency injection, and invoke `SetupRoutes` in handler layer.
2. Define interfaces for each domain layer in `internal/domain`.
3. Add swagger annotations to each handler `router.go` for API documentation.
4. Review the updated memory bank files (`architecture.md`, `tech.md`, `context.md`) with the user for accuracy and approval.
5. Proceed with regenerating the SQLC code to apply the new structure.
6. Migrate existing identifiers and default primary key generation to use UUIDv7 for index-friendly primary keys.
7. Integrate gofiber/swagger to generate and serve API documentation.
8. Configure and use samber/do for application dependency injection.