# Current Context

## Project State
Go Fiber Skeleton template is in **active development** with core infrastructure and two example domains (user, post) implemented. The project has a working foundation but is not yet feature-complete per the brief requirements.

## Recent Implementation Status

### ✅ Completed Features
1. **Core Infrastructure**
   - Fiber v2 web server with graceful shutdown
   - PostgreSQL connection via database/sql + pgx driver
   - Zerolog structured logging with request ID tracking
   - JSend response format standardization
   - Docker Compose environment (PostgreSQL 18 + Valkey 8)
   - Multi-stage Dockerfile with Alpine-based production image

2. **Configuration System**
   - godotenv-based config loading (.env file support)
   - Environment variable override capability
   - Config struct with Port, DatabaseDSN, JWTSecret

3. **Database Tooling**
   - sqlc integration with complete configuration (sqlc.yaml)
   - Four migration files created (users, roles, threads/posts/comments, sessions)
   - Type-safe query generation for users and posts
   - Makefile command for sqlc generation

4. **Authentication Domain (User)**
   - Complete 4-layer implementation (entity, repository, usecase, handler)
   - User registration with bcrypt password hashing
   - JWT-based login with 72-hour token expiry
   - Auth middleware for protected routes
   - Swagger documentation annotations

5. **Post Domain**
   - Complete 4-layer implementation demonstrating Clean Architecture
   - CRUD operations with ownership validation
   - Mock generation setup (//go:generate annotations)
   - Mock files generated for repository and usecase

6. **Middleware**
   - Request ID generation
   - Structured logging with duration metrics
   - API rate limiting
   - Auth-specific rate limiting
   - Panic recovery

7. **API Documentation**
   - Swagger/OpenAPI integration via swaggo
   - Auto-generation from code comments
   - Swagger UI accessible at /swagger

8. **Development Tooling**
   - Comprehensive Makefile (fmt, build, run, test, sqlc, lint, ci)
   - Docker Compose with health checks
   - Test infrastructure (placeholder tests exist)

### ⚠️ Partially Implemented
1. **Testing Strategy**
   - Mock generation infrastructure ready (mockgen annotations present)
   - Mock files generated for post domain
   - Placeholder test file exists (`internal/post/tests/post_test.go`)
   - **Missing**: Actual unit tests with mock usage
   - **Missing**: Repository tests with go-sqlmock (dependency not added)

2. **Database Migrations**
   - Migration files created
   - **Missing**: golang-migrate integration
   - **Missing**: Makefile migrate command implementation (currently placeholder)

3. **Dependency Injection**
   - Brief specifies Uber fx
   - **Current**: Manual DI in router.go (works but not using fx)

### ❌ Not Yet Implemented
1. **Hot Reloading**
   - Air mentioned in brief but not configured
   - No .air.toml configuration file

2. **Viper Configuration**
   - Brief specifies Viper for config management
   - **Current**: Using godotenv instead (simpler but less flexible)

3. **Redis/Valkey Integration**
   - Infrastructure running in Docker Compose
   - **Missing**: Go client library and integration code

4. **Comprehensive Testing**
   - go-sqlmock not in dependencies
   - No actual test implementations beyond placeholder

5. **Additional Domains**
   - Migrations exist for roles, threads, comments
   - **Missing**: Implementation code for these domains

6. **Domain Addition Guide**
   - Brief requires "clear instructions" for adding domains
   - **Missing**: Documentation in README or separate guide

## Active Work Focus
**Status**: Foundation complete, awaiting direction for next phase

### Immediate Next Steps (Not Yet Started)
1. Implement comprehensive test suite using existing mock infrastructure
2. Add golang-migrate integration for database migrations
3. Document the process for adding new domains
4. Implement remaining domains (threads, comments, roles)
5. Add Air for hot-reloading in development

## Current File State

### Recently Modified/Generated
- `internal/db/*.go` - sqlc-generated code (current with queries)
- `internal/post/repository/mocks/repository.go` - Generated mocks
- `internal/post/usecase/mocks/usecase.go` - Generated mocks
- `docs/*` - Swagger documentation (generated)

### Needs Attention
- `Makefile` - migrate command is placeholder, needs implementation
- `README.md` - Basic, needs domain addition guide
- `internal/post/tests/post_test.go` - Placeholder, needs real tests
- `go.mod` - Missing: go-sqlmock, Air, potentially golang-migrate

## Known Gaps Between Brief and Implementation

### Critical Gaps
1. **Testing Not Demonstrated**: Brief emphasizes "comprehensive testing strategy" with examples, current has only infrastructure
2. **Migration Tool Not Integrated**: Brief specifies golang-migrate, not connected
3. **No Domain Addition Guide**: Brief requires "clear instructions," not documented
4. **Viper vs godotenv**: Brief specifies Viper, implementation uses simpler godotenv

### Medium Priority Gaps
1. **Uber fx**: Mentioned in brief, not used (manual DI works but not as specified)
2. **Air**: Development tool mentioned, not configured
3. **Incomplete Domains**: Roles, threads, comments have migrations but no code

### Low Priority Gaps
1. **Redis Integration**: Infrastructure ready but client not connected
2. **Additional Middleware**: Could expand rate limiting, add CORS, etc.

## Technical Decisions Made

### Deviations from Brief
1. **godotenv instead of Viper**: Simpler approach chosen, works for current needs
2. **Manual DI instead of Uber fx**: Less complexity, easier to understand for template users
3. **Minimal initial domains**: Focus on quality examples (user, post) rather than full forum

### Architectural Patterns Established
1. **4-Layer Domain Structure**: entity → repository → usecase → handler (strictly enforced)
2. **Interface-Driven Design**: All repositories and usecases defined as interfaces
3. **Mock Generation Pattern**: `//go:generate mockgen` annotations on interfaces
4. **JSend Response Standard**: All handlers use response.Success/Fail/Error
5. **UUID v7**: Time-sortable UUIDs for all entities

## Blockers & Questions

### Current Blockers: None
Project is functional and deployable as-is. Gaps are features awaiting implementation.

### Pending Decisions
1. Should Viper replace godotenv to match brief exactly?
2. Should Uber fx be added for DI, or keep manual approach for simplicity?
3. Priority order for implementing remaining domains vs. completing testing suite?
4. Should Air configuration be added, or leave to individual developer preference?

## Environment Notes

### Development Setup Requirements
1. Docker and Docker Compose installed
2. Go 1.24.6+ installed
3. sqlc installed for regenerating queries
4. swag installed for regenerating API docs
5. mockgen installed for regenerating mocks

### Current Database Schema
- **users**: id, username, email, password_hash, created_at, updated_at
- **roles**: (migration exists, table structure unknown without reading file)
- **threads**: (migration exists, structure partially known - has posts relation)
- **posts**: id, thread_id, user_id, content, created_at, updated_at
- **comments**: (migration exists, table structure unknown)
- **sessions**: (migration exists, table structure unknown)

## Integration Points

### External Services (via Docker Compose)
- PostgreSQL: localhost:5432 (db service)
- Valkey/Redis: localhost:6379 (redis service, not yet used by app)

### API Endpoints Implemented
- `GET /health` - Liveness check
- `GET /ready` - Readiness check (DB ping)
- `GET /swagger/*` - API documentation
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/posts` - Create post (protected)
- `GET /api/v1/posts/:id` - Get post
- `GET /api/v1/posts/user/:user_id` - List user's posts
- `PUT /api/v1/posts/:id` - Update post (protected, ownership check)
- `DELETE /api/v1/posts/:id` - Delete post (protected, ownership check)

### Stub Routes (Not Implemented)
- `GET /api/v1/users`
- `GET /api/v1/threads`
- `GET /api/v1/comments`

## Next Session Priorities

**When work resumes, prioritize:**
1. Review and validate Memory Bank files with user
2. Decide on Viper vs godotenv and Uber fx vs manual DI
3. Implement test suite with mocks (demonstrates testing strategy from brief)
4. Integrate golang-migrate for proper migration management
5. Document domain addition process in README or separate guide