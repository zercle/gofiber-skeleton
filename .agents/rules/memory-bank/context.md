# Current Context: Go Fiber Microservice Template

## Development Phase
**Phase:** 0-to-1 Development - Specification & Template Scaffolding
**Status:** Specification Complete, Implementation Pending
**Last Updated:** 2025-10-24

## Project Status Summary
The template project is in the planning and specification phase with comprehensive documentation already established. All architectural decisions have been made through the brief.md specification. The next phase is implementation of the core template structure and example components.

**Current State:**
- ✅ Complete specification (brief.md) - 122 lines
- ✅ Project structure defined in detail
- ✅ Technology stack finalized
- ✅ Architectural patterns documented
- ✅ Code quality standards established
- ❌ No source code implementation yet
- ❌ No example components (handler, usecase, repository)
- ❌ No configuration management setup
- ❌ No database layer (migrations, queries)
- ❌ No API documentation (Swagger)
- ❌ No build automation (Makefile)
- ❌ No Docker support
- ❌ No tests

## Recent Changes
**Date: 2025-10-24**
- Memory Bank initialization (product.md, architecture.md, tech.md, context.md, specs.md created)
- Analyzed full codebase structure and documented in memory bank
- Confirmed specification completeness in brief.md

**Previous Commits:**
- 2025-10-23: Updated AGENTS.md with optimized agentic AI coder ruleset
- 2025-10-23: Squashed commits documenting transition from Echo-based to Fiber-based template

## Key Architectural Decisions

### 1. Framework Choice: Fiber
**Decision:** Use Fiber framework for HTTP layer
**Rationale:** Fastest Go web framework, perfect for microservices, inspired by Express.js (familiar to many developers)
**Trade-off:** Smaller ecosystem than some alternatives, but sufficient for microservices

### 2. Layered Architecture (Handler → Usecase → Repository)
**Decision:** Use 4-layer architecture pattern
**Rationale:** Clear separation of concerns, easy to test, aligns with domain-driven design concepts
**Layers:**
- Handler: HTTP concerns only (routing, request/response mapping)
- Usecase: Business logic and orchestration
- Repository: Data access abstraction
- Domain: Core domain models

### 3. Dependency Injection: samber/do v2
**Decision:** Use samber/do for DI container
**Rationale:** Lightweight, type-safe, zero-runtime-cost, perfect for simple services
**Alternative Considered:** Wire (better for complex apps, but overkill for template)

### 4. Database: sqlc
**Decision:** Type-safe SQL code generation
**Rationale:** Prevents SQL injection, compile-time error checking, no ORM learning curve
**Implication:** Developers write SQL, sqlc generates Go code

### 5. Authentication: JWT Primary, OIDC Alternative
**Decision:** JWT for primary auth, OIDC pattern available
**Rationale:** JWT is stateless (fits microservices), OIDC for enterprise requirements
**Security:** Leverages golang-jwt/jwt library with secure token handling

### 6. Testing: Interface-driven with uber-go/mock
**Decision:** All interfaces auto-generate mocks via //go:generate
**Rationale:** Eliminates mock maintenance, ensures testability from design time
**Database Testing:** go-sqlmock for repository testing

### 7. Logging: Structured Logging with slog
**Decision:** Use slog with context propagation
**Rationale:** Built-in Go library, structured logging industry standard, supports context
**Implication:** All logs are machine-readable JSON

### 8. Configuration: Environment Variables via Viper
**Decision:** Environment-variable-driven with Viper loading
**Rationale:** 12-factor app principles, simple validation, deployment-friendly
**Pattern:** Load → Validate → Inject via DI container

## Next Steps (Prioritized)

### Phase 1: Core Infrastructure (Weeks 1-2)
1. **Setup Go Project Structure**
   - Create cmd/api/main.go entry point
   - Initialize go.mod with core dependencies
   - Setup graceful shutdown pattern
   - Configure slog structured logging
   - Create domain/ package with core interfaces

2. **Configuration Management**
   - Implement config/config.go with Viper loading
   - Create .env.example template
   - Add config validation logic
   - Setup environment-based initialization

3. **Build Automation**
   - Create comprehensive Makefile
   - Setup test command with coverage
   - Add lint command (golangci-lint v2)
   - Add docker build command
   - Add migrate command for DB migrations

### Phase 2: Layer Implementations (Weeks 2-3)
1. **Handler Layer Example**
   - Create handler/user.go with CRUD endpoint examples
   - Add Swagger/OpenAPI annotations
   - Demonstrate JSend response format
   - Setup input validation patterns

2. **Usecase Layer Example**
   - Create usecase/user.go with business logic
   - Demonstrate transaction handling
   - Show dependency injection usage
   - Add error handling patterns

3. **Repository Layer Example**
   - Create repository/user.go with database access
   - Setup sqlc query generation
   - Demonstrate transaction support
   - Add database mocking pattern

### Phase 3: Database & API Documentation (Week 3)
1. **Database Setup**
   - Create sql/migrations/ directory structure
   - Create example migration (users table)
   - Setup sql/queries/ for sqlc
   - Configure golang-migrate integration

2. **API Documentation**
   - Generate Swagger UI from handler annotations
   - Setup /docs endpoint
   - Create docs/swagger.yaml
   - Add example requests/responses

### Phase 4: Testing & Deployment (Week 4)
1. **Testing Framework**
   - Setup test directory structure
   - Create mock generation examples
   - Add integration test examples
   - Setup test fixtures for database

2. **Containerization**
   - Create Dockerfile with multi-stage build
   - Setup Docker Compose for dev environment
   - Add .dockerignore
   - Demonstrate graceful shutdown in container

3. **CI/CD Pipeline**
   - Create GitHub Actions workflow
   - Setup automated testing
   - Add linting checks
   - Add security scanning

### Phase 5: Documentation & Examples (Week 5)
1. **Developer Documentation**
   - Create comprehensive README.md
   - Add ARCHITECTURE.md with detailed design
   - Create GETTING_STARTED.md guide
   - Add CONTRIBUTING.md guidelines

2. **Example Implementations**
   - Complete user management example
   - Demonstrate common patterns (auth, errors, validation)
   - Add inline code comments explaining decisions
   - Create troubleshooting guide

3. **ADR Documentation**
   - Document all key architectural decisions
   - Record rationale and trade-offs
   - List alternatives considered
   - Create ADR index

## Known Issues & Blockers
**None at present** - Project is in specification phase with clear path forward.

## Technical Debt
**None yet** - Project hasn't started implementation.

## Team Knowledge & Context
- **Architecture Understanding:** Complete (documented in brief.md)
- **Technology Stack:** Finalized (documented in tech.md)
- **Code Examples:** To be created (Phase 2)
- **Testing Patterns:** To be established (Phase 4)

## Memory Bank Maintenance
- **Last Comprehensive Update:** 2025-10-24 (initial creation)
- **Files to Review Before Next Major Change:** brief.md, tech.md, architecture.md
- **Sync Points:** After each phase completion, update context.md with progress

## Quick Reference Links
- **Specification:** brief.md (project requirements)
- **Product Definition:** product.md (user goals and success criteria)
- **Architecture Details:** architecture.md (component design and ADRs)
- **Technology Details:** tech.md (stack, dependencies, setup)
- **Feature Specifications:** specs.md (detailed feature requirements)

## Questions for Clarification
1. **Example Domain:** Should we use "User" as the example domain or something else (e.g., "Task", "Item")?
2. **Database Choice:** Start with MariaDB example or PostgreSQL? (Brief mentions both)
3. **Priority:** Focus on documentation/examples or getting core structure working first?
4. **API Scope:** Should example include any advanced features (pagination, filters) or keep minimal?
5. **Testing:** Start with unit tests or integration tests as primary focus?
