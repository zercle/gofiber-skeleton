# Common Task Workflows: Go Fiber Forum Backend

This document outlines the workflow for reinitializing the Go Fiber backend template to a feature-based structure and implementing a simple blog post feature with authentication.

## 0. Reinitialize Project Structure

**Task Name:** `project_reinit_feature_based`

**Description:** Migrate codebase to feature-based folder structure and implement simple blog post feature with authentication.

**Files to Update/Create:**
- `architecture.md`
- `context.md`
- `tasks.md`
- `features.md`
- Project directories: `cmd/app/main.go`, `internal/user/`, `internal/post/`, `db/migrations/`, `db/queries/`

**Step-by-step Process:**

1. Create feature directories under `internal/`: `user`, `post` with subdirectories: `handler/`, `usecase/`, `repository/`, `entity/`, `model/`, `tests/`.
2. Move existing domain code into the new feature directories and adjust import paths.
3. Add migration SQL scripts to `db/migrations` for `users` and `posts` tables.
4. Add SQL query definitions to `db/queries` and run `sqlc` to generate code.
5. Implement user authentication: registration and login handlers with Bcrypt and JWT.
6. Implement post handlers: CRUD operations secured with authentication middleware.
7. Update application entrypoint in `cmd/app/main.go` to register feature routes.
8. Update `architecture.md` with new feature-based architecture overview.
9. Update `context.md` with reinitialization context and goals.
10. Create `features.md` detailing feature-based folder organization.
11. Overwrite `tasks.md` with this workflow for future reference.
12. Verify code quality:
```bash
go generate ./... && golangci-lint run --fix ./... && go clean -testcache && go test -v -race ./...