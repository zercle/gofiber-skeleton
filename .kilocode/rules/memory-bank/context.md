# Context

**Current Work Focus**
Finalizing the architectural update to centralize SQLC-generated code and query files.

**Recent Changes**
- Centralized all SQLC-generated code into the `internal/infrastructure/sqlc` directory.
- Moved all SQL query files to a single `queries` directory at the project root.
- Updated `sqlc.yaml` to reflect the new centralized structure.
- Updated `architecture.md` and `tech.md` to document the new SQLC configuration and file locations.

**Next Steps**
1. Review the updated memory bank files (`architecture.md`, `tech.md`, `context.md`) with the user for accuracy and approval.
2. Proceed with regenerating the SQLC code to apply the new structure.
3. Migrate existing identifiers and default primary key generation to use UUIDv7 for index-friendly primary keys.