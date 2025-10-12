# **Tasks & Workflows Documentation: Go Fiber Skeleton**

## **1. Development Setup Tasks**

### **Task 1: Initialize New Project from Template**
**Purpose**: Set up a new project using this template foundation
**Duration**: 5-10 minutes
**Dependencies**: Git, Go 1.25.0+, Docker

**Step-by-Step Procedure**:
1. **Clone Repository**
   ```bash
   git clone https://github.com/zercle/gofiber-skeleton.git my-project
   cd my-project
   ```

2. **Update Module Information**
   - Edit `go.mod` to change module name
   - Update `README.md` with project-specific information
   - Modify `.env.example` with project-specific configuration

3. **Initialize Development Environment**
   ```bash
   make dev
   ```

4. **Verify Setup**
   - Check server starts at `http://localhost:3000`
   - Verify API docs at `http://localhost:3000/swagger/`
   - Run tests: `make test`

**Affected Files**: `go.mod`, `README.md`, `.env.example`
**Success Criteria**: Development server running, tests passing, API documentation accessible

---

### **Task 2: Add New Business Domain**
**Purpose**: Create a new domain following the established patterns
**Duration**: 30-60 minutes
**Dependencies**: Existing project structure, database connection

**Step-by-Step Procedure**:
1. **Create Domain Structure**
   ```bash
   mkdir -p internal/domains/{domain_name}/entity
   mkdir -p internal/domains/{domain_name}/repository
   mkdir -p internal/domains/{domain_name}/usecase
   mkdir -p internal/domains/{domain_name}/handler
   ```

2. **Define Entity**
   - Create entity file: `internal/domains/{domain_name}/entity/{entity}.go`
   - Define struct with validation tags
   - Add business logic methods

3. **Create Repository Interface**
   - Create repository interface: `internal/domains/{domain_name}/repository/interface.go`
   - Define CRUD operations
   - Add custom query methods
   - Add `//go:generate mockgen` annotation

4. **Implement Usecases**
   - Create usecase file: `internal/domains/{domain_name}/usecase/{usecase}.go`
   - Implement business logic
   - Add validation and error handling
   - Write comprehensive tests

5. **Create HTTP Handlers**
   - Create handler file: `internal/domains/{domain_name}/handler/{handler}.go`
   - Implement HTTP endpoints
   - Add Swagger documentation
   - Add request/response validation

6. **Add Database Integration**
   - Create SQL migration: `db/migrations/{version}_{description}.up.sql`
   - Create SQL queries: `db/queries/{domain_name}.sql`
   - Generate Go code: `make generate`

7. **Wire Dependencies**
   - Add domain to DI container
   - Register routes in router
   - Update configuration if needed

**Affected Files**:
- `internal/domains/{domain_name}/**/*`
- `db/migrations/*`
- `db/queries/*`
- Router configuration
- DI container setup

**Success Criteria**:
- Domain CRUD operations working
- API endpoints functional and documented
- Tests passing with 90%+ coverage
- Integration with existing architecture seamless

---

### **Task 3: Database Migration Management**
**Purpose**: Create and apply database schema changes
**Duration**: 10-30 minutes
**Dependencies**: PostgreSQL connection, migrate tool

**Step-by-Step Procedure**:
1. **Create Migration File**
   ```bash
   make create-migration NAME=add_users_table
   ```

2. **Write Migration SQL**
   - Edit up migration file with schema changes
   - Edit down migration file with rollback logic
   - Test SQL in development database

3. **Apply Migration**
   ```bash
   make migrate-up
   ```

4. **Verify Migration**
   - Check database schema
   - Test new schema with application
   - Verify rollback works: `make migrate-down 1`

5. **Update sqlc Queries**
   - Add SQL queries to `db/queries/`
   - Generate Go code: `make generate`
   - Update repository implementations

**Affected Files**: `db/migrations/*`, `db/queries/*`, repository implementations
**Success Criteria**: Migration applied successfully, schema verified, queries working

---

## **2. Development Workflow Tasks**

### **Task 4: Run Local Development Environment**
**Purpose**: Start development environment with all services
**Duration**: 2-3 minutes
**Dependencies**: Docker, Docker Compose

**Step-by-Step Procedure**:
1. **Start Services**
   ```bash
   make dev
   ```

2. **Verify Services**
   - Application: http://localhost:3000
   - Database: localhost:5432
   - Cache: localhost:6379
   - API Docs: http://localhost:3000/swagger/

3. **Monitor Logs**
   ```bash
   docker compose logs -f
   ```

4. **Stop Services**
   ```bash
   make stop
   ```

**Affected Files**: `compose.yml`, Docker configuration
**Success Criteria**: All services running, logs accessible, endpoints responding

---

### **Task 5: Run Tests and Coverage**
**Purpose**: Execute test suite and verify coverage
**Duration**: 2-5 minutes
**Dependencies**: Test dependencies, database

**Step-by-Step Procedure**:
1. **Run All Tests**
   ```bash
   make test
   ```

2. **Check Coverage**
   ```bash
   make coverage
   ```

3. **View Coverage Report**
   ```bash
   make coverage-html
   # Open coverage.html in browser
   ```

4. **Run Specific Test**
   ```bash
   go test ./internal/domains/user/usecase/...
   ```

5. **Run Tests with Verbose Output**
   ```bash
   make test-verbose
   ```

**Affected Files**: Test files, coverage reports
**Success Criteria**: All tests passing, coverage >90%, report generated

---

### **Task 6: Code Quality and Linting**
**Purpose**: Ensure code quality standards
**Duration**: 1-2 minutes
**Dependencies**: golangci-lint, development tools

**Step-by-Step Procedure**:
1. **Format Code**
   ```bash
   make fmt
   ```

2. **Run Linter**
   ```bash
   make lint
   ```

3. **Fix Linting Issues**
   - Review linting output
   - Fix code formatting and style issues
   - Re-run linter to verify fixes

4. **Run Security Scan**
   ```bash
   make security
   ```

5. **Run Tidy**
   ```bash
   make tidy
   ```

**Affected Files**: All Go source files
**Success Criteria**: No linting errors, security scan clean, dependencies tidy

---

## **3. Build and Deployment Tasks**

### **Task 7: Build Application**
**Purpose**: Build application for deployment
**Duration**: 1-2 minutes
**Dependencies**: Go build tools

**Step-by-Step Procedure**:
1. **Build for Development**
   ```bash
   make build
   ```

2. **Build for Production**
   ```bash
   make build-prod
   ```

3. **Build for Multiple Platforms**
   ```bash
   make build-all
   ```

4. **Verify Build**
   - Check binary exists
   - Test binary execution
   - Verify version information

**Affected Files**: Binary files, build artifacts
**Success Criteria**: Binary created, executable, version info correct

---

### **Task 8: Docker Deployment**
**Purpose**: Build and deploy Docker containers
**Duration**: 3-5 minutes
**Dependencies**: Docker, Docker Hub access

**Step-by-Step Procedure**:
1. **Build Docker Image**
   ```bash
   make docker-build
   ```

2. **Tag Image for Registry**
   ```bash
   make docker-tag
   ```

3. **Push to Registry**
   ```bash
   make docker-push
   ```

4. **Deploy Container**
   ```bash
   docker compose -f compose.prod.yml up -d
   ```

5. **Verify Deployment**
   - Check container status
   - Test health endpoints
   - Monitor logs

**Affected Files**: Docker image, container registry, production compose file
**Success Criteria**: Container built, pushed, deployed, and running healthy

---

## **4. Maintenance Tasks**

### **Task 9: Update Dependencies**
**Purpose**: Keep dependencies up to date
**Duration**: 5-10 minutes
**Dependencies**: Go modules, internet access

**Step-by-Step Procedure**:
1. **Check for Updates**
   ```bash
   go list -u -m all
   ```

2. **Update Dependencies**
   ```bash
   make update-deps
   ```

3. **Test Updates**
   ```bash
   make test
   make build
   ```

4. **Review Changes**
   - Check for breaking changes
   - Review release notes
   - Update documentation if needed

5. **Commit Updates**
   ```bash
   git add go.mod go.sum
   git commit -m "chore: update dependencies"
   ```

**Affected Files**: `go.mod`, `go.sum`, documentation
**Success Criteria**: Dependencies updated, tests passing, build successful

---

### **Task 10: Database Maintenance**
**Purpose**: Maintain database health and performance
**Duration**: 15-30 minutes
**Dependencies**: Database access, monitoring tools

**Step-by-Step Procedure**:
1. **Check Database Status**
   ```bash
   make db-status
   ```

2. **Run Database Health Check**
   ```bash
   make db-health
   ```

3. **Backup Database**
   ```bash
   make db-backup
   ```

4. **Analyze Query Performance**
   ```bash
   make db-analyze
   ```

5. **Optimize Database**
   ```bash
   make db-optimize
   ```

6. **Review Migration Status**
   ```bash
   make db-status
   ```

**Affected Files**: Database schema, backup files
**Success Criteria**: Database healthy, backup created, performance optimized

---

## **5. Troubleshooting Tasks**

### **Task 11: Debug Application Issues**
**Purpose**: Diagnose and fix application problems
**Duration**: 10-30 minutes
**Dependencies**: Application logs, debugging tools

**Step-by-Step Procedure**:
1. **Check Application Logs**
   ```bash
   make logs
   ```

2. **Check Service Health**
   ```bash
   make health
   ```

3. **Run in Debug Mode**
   ```bash
   DEBUG=true make dev
   ```

4. **Test Individual Components**
   - Test database connection
   - Test cache connection
   - Test API endpoints individually

5. **Enable Verbose Logging**
   ```bash
   LOG_LEVEL=debug make dev
   ```

6. **Review Configuration**
   - Check environment variables
   - Verify configuration files
   - Validate secret values

**Affected Files**: Application logs, configuration files
**Success Criteria**: Issue identified, root cause found, fix implemented

---

### **Task 12: Performance Issues**
**Purpose**: Diagnose and resolve performance problems
**Duration**: 20-60 minutes
**Dependencies**: Monitoring tools, load testing

**Step-by-Step Procedure**:
1. **Monitor Resource Usage**
   ```bash
   make monitor
   ```

2. **Run Performance Tests**
   ```bash
   make load-test
   ```

3. **Profile Application**
   ```bash
   make profile
   ```

4. **Analyze Database Queries**
   ```bash
   make db-profile
   ```

5. **Check Memory Usage**
   ```bash
   make memory-profile
   ```

6. **Review Performance Metrics**
   - Response times
   - Memory usage
   - CPU utilization
   - Database performance

**Affected Files**: Performance reports, profiling data
**Success Criteria**: Bottleneck identified, optimization implemented, performance improved

---

## **6. Documentation Tasks**

### **Task 13: Update API Documentation**
**Purpose**: Keep API documentation current
**Duration**: 5-10 minutes
**Dependencies**: Swagger annotations, swag tool

**Step-by-Step Procedure**:
1. **Add Swagger Comments**
   - Add `@Summary` annotations
   - Add `@Description` annotations
   - Add `@Tags` annotations
   - Add `@Accept` and `@Produce` annotations
   - Add `@Param` and `@Success` annotations

2. **Generate Documentation**
   ```bash
   make docs
   ```

3. **Review Generated Docs**
   - Check documentation at `/swagger/`
   - Verify all endpoints documented
   - Validate examples and schemas

4. **Test API Examples**
   - Try examples in Swagger UI
   - Verify responses match documentation
   - Update examples if needed

**Affected Files**: Handler files, generated documentation
**Success Criteria**: Documentation generated, complete, accurate examples

---

### **Task 14: Memory Bank Maintenance**
**Purpose**: Keep Memory Bank files current and accurate
**Duration**: 10-20 minutes
**Dependencies**: Project analysis, documentation updates

**Step-by-Step Procedure**:
1. **Analyze Current State**
   ```bash
   update memory bank
   ```

2. **Review Memory Bank Content**
   - Check accuracy of `architecture.md`
   - Review current state in `context.md`
   - Validate technology stack in `tech.md`
   - Review product features in `product.md`

3. **Update Context Information**
   ```bash
   refresh context
   ```

4. **Add New Tasks**
   ```bash
   add task: [Task Name]
   ```

5. **Validate Memory Bank**
   - Check file consistency
   - Validate cross-references
   - Review completeness

**Affected Files**: Memory Bank files
**Success Criteria**: Memory Bank current, accurate, complete

---

## **7. Automation and Scripts**

### **Task 15: Create Custom Development Scripts**
**Purpose**: Add project-specific automation
**Duration**: 30-60 minutes
**Dependencies**: Shell scripting, Makefile knowledge

**Step-by-Step Procedure**:
1. **Identify Repetitive Tasks**
   - Common development workflows
   - Frequent manual processes
   - Complex multi-step operations

2. **Create Shell Scripts**
   - Add scripts to `scripts/` directory
   - Make scripts executable
   - Add error handling
   - Include logging and feedback

3. **Integrate with Makefile**
   - Add new Make targets
   - Document usage
   - Chain with existing targets
   - Add dependencies between targets

4. **Test Scripts**
   - Test in development environment
   - Verify error handling
   - Test edge cases
   - Document expected behavior

5. **Add Documentation**
   - Update README.md
   - Add usage examples
   - Document prerequisites
   - Provide troubleshooting guide

**Affected Files**: `scripts/`, `Makefile`, `README.md`
**Success Criteria**: Scripts functional, integrated, documented

---

## **Task Templates**

### **New Domain Template**
Use this template when adding new business domains:

```markdown
### **Task: Add [Domain Name] Domain**
**Files to Create**:
- `internal/domains/[domain]/entity/[entity].go`
- `internal/domains/[domain]/repository/interface.go`
- `internal/domains/[domain]/usecase/[usecase].go`
- `internal/domains/[domain]/handler/[handler].go`

**Database Components**:
- Migration file: `db/migrations/[version]_[description].sql`
- Query file: `db/queries/[domain].sql`

**Integration Points**:
- Add to DI container in `internal/app/app.go`
- Add routes in `internal/router/router.go`
- Update configuration if needed

**Testing Requirements**:
- Unit tests for usecases
- Integration tests for handlers
- Mock generation with `//go:generate`
- Coverage target: 90%+
```

### **Maintenance Task Template**
Use this template for regular maintenance:

```markdown
### **Task: [Maintenance Activity]**
**Frequency**: [Daily/Weekly/Monthly]
**Prerequisites**: [Required tools/permissions]
**Expected Duration**: [Time estimate]
**Success Criteria**: [What constitutes completion]

**Steps**:
1. [Step 1]
2. [Step 2]
3. [Step 3]

**Rollback Plan**: [How to undo if needed]
**Verification**: [How to verify success]
```

---

This tasks documentation provides comprehensive workflows for common development activities. Tasks are designed to be repeatable, well-documented, and include clear success criteria and verification steps.