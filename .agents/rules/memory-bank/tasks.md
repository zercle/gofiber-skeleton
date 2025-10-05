# Development Tasks & Workflows

## Common Development Workflows

### 1. Adding a New Domain

**Objective**: Add a new business domain following Clean Architecture patterns

**Prerequisites**:
- Development tools installed (`make install-tools`)
- Understanding of Clean Architecture principles
- Database schema design for the domain

**Steps**:
1. **Create Directory Structure**
   ```bash
   mkdir -p internal/{domain}/entity
   mkdir -p internal/{domain}/repository/mocks
   mkdir -p internal/{domain}/usecase/mocks
   mkdir -p internal/{domain}/handler
   mkdir -p internal/{domain}/tests
   ```

2. **Define Entity**
   - Create `internal/{domain}/entity/{model}.go`
   - Define domain model with UUID primary keys
   - Include audit fields (created_at, updated_at)

3. **Create Database Migration**
   ```bash
   make migrate-create name=create_{table_name}
   ```
   - Edit generated migration files in `db/migrations/`
   - Include proper indexes and constraints

4. **Define SQL Queries**
   - Create `db/queries/{domain}.sql`
   - Follow sqlc naming conventions
   - Include CRUD operations and business-specific queries

5. **Generate Type-Safe Code**
   ```bash
   make sqlc
   ```

6. **Implement Repository**
   - Create `internal/{domain}/repository/postgres.go`
   - Define interface first, then implementation
   - Add `//go:generate` annotation for mocks

7. **Implement Usecase**
   - Create `internal/{domain}/usecase/{domain}.go`
   - Define interface with business logic methods
   - Implement business rules and validation
   - Add `//go:generate` annotation for mocks

8. **Generate Mocks**
   ```bash
   make generate-mocks
   ```

9. **Implement HTTP Handlers**
   - Create `internal/{domain}/handler/{handler}.go`
   - Add Swagger annotations for documentation
   - Implement request validation
   - Use JSend response format

10. **Register Routes**
    - Edit `internal/server/router.go`
    - Add domain imports and initialization
    - Register routes with appropriate middleware

11. **Write Tests**
    - Create unit tests for usecases with mocks
    - Test repository logic with sqlmock
    - Add integration tests for handlers

12. **Finalize**
    ```bash
    make migrate-up
    make test
    make generate-docs
    ```

**Time Estimate**: 15-20 minutes for simple CRUD domain

**Common Pitfalls**:
- Forgetting foreign key constraints in migrations
- Missing validation tags on request structs
- Not regenerating mocks after interface changes
- Circular dependencies between domains

---

### 2. Database Schema Changes

**Objective**: Modify database schema safely with proper migrations

**Types of Changes**:
- Adding new tables
- Adding columns to existing tables
- Modifying constraints
- Creating indexes
- Data migrations

**Steps**:
1. **Create Migration**
   ```bash
   make migrate-create name=descriptive_change_name
   ```

2. **Write Up Migration**
   - Add DDL statements in `XXX_descriptive_change_name.up.sql`
   - Include proper error handling
   - Consider performance impact

3. **Write Down Migration**
   - Add rollback statements in `XXX_descriptive_change_name.down.sql`
   - Ensure rollback is safe and complete

4. **Test Migration**
   ```bash
   # Test on local development database
   make migrate-up
   
   # Test rollback
   make migrate-down
   make migrate-up
   ```

5. **Update SQL Queries**
   - Modify relevant `.sql` files in `db/queries/`
   - Add new queries for new columns/tables

6. **Regenerate Code**
   ```bash
   make sqlc
   ```

7. **Update Repository Layer**
   - Modify repository implementations
   - Update entity mappings
   - Handle new fields in usecases

8. **Test Changes**
   ```bash
   make test
   make migrate-version
   ```

**Best Practices**:
- Always write down migrations
- Test migrations on copy of production data
- Consider backward compatibility
- Use transactions for complex changes

---

### 3. Adding New API Endpoints

**Objective**: Add new HTTP endpoints to existing domain

**Steps**:
1. **Define Request/Response DTOs**
   - Create request structs with validation tags
   - Define response structures
   - Add Swagger annotations

2. **Implement Handler Method**
   - Add method to handler struct
   - Include proper error handling
   - Use JSend response format
   - Add authentication/authorization if needed

3. **Add Usecase Method**
   - Implement business logic in usecase
   - Add input validation
   - Handle business rules
   - Return appropriate errors

4. **Update Repository (if needed)**
   - Add new SQL queries
   - Regenerate sqlc code
   - Implement repository methods

5. **Register Route**
   - Add route in router registration function
   - Apply appropriate middleware
   - Consider rate limiting

6. **Generate Documentation**
   ```bash
   make generate-docs
   ```

7. **Write Tests**
   - Unit tests for usecase logic
   - Handler tests with mock usecase
   - Integration tests for full flow

**Validation Requirements**:
- All request DTOs must have validation tags
- Authentication required for protected endpoints
- Authorization checks for resource access
- Proper HTTP status codes

---

### 4. Testing Workflows

### Unit Testing

**Objective**: Test business logic in isolation

**Pattern**:
```go
func TestUsecase_Method(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mock.NewMockRepository(ctrl)
    usecase := NewUsecase(mockRepo)
    
    // Setup expectations
    mockRepo.EXPECT().
        Method(gomock.Any(), gomock.Any()).
        Return(expectedResult, nil).
        Times(1)
    
    // Execute test
    result, err := usecase.Method(ctx, input)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

**Requirements**:
- Test all public methods
- Test error scenarios
- Test edge cases
- Use table-driven tests for multiple scenarios

### Integration Testing

**Objective**: Test full request/response flow

**Setup**:
1. Use test database with migrations
2. Real HTTP client for endpoint testing
3. Test fixtures for consistent data

**Pattern**:
```go
func TestHandler_Endpoint(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Setup test server
    app := setupTestApp(db)
    
    // Make request
    req := httptest.NewRequest("POST", "/api/v1/endpoint", body)
    resp, err := app.Test(req)
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}
```

---

### 5. Performance Optimization Tasks

### Database Optimization

**Query Optimization**:
1. Analyze slow queries with `EXPLAIN ANALYZE`
2. Add appropriate indexes
3. Optimize JOIN operations
4. Use connection pooling effectively

**Migration Performance**:
1. Test migrations on large datasets
2. Use batch operations for data migrations
3. Consider zero-downtime migrations for production

### Application Optimization

**Memory Management**:
1. Profile memory usage with `pprof`
2. Optimize goroutine usage
3. Implement object pooling where appropriate
4. Monitor garbage collection

**Concurrency Optimization**:
1. Use worker pools for CPU-bound tasks
2. Implement proper channel buffering
3. Avoid goroutine leaks
4. Use context for cancellation

---

### 6. Security Enhancement Tasks

### Authentication & Authorization

**JWT Security**:
1. Implement token refresh mechanism
2. Add token blacklisting for logout
3. Use short-lived tokens with refresh
4. Implement proper secret rotation

**Authorization**:
1. Implement role-based access control (RBAC)
2. Add resource-level permissions
3. Implement ownership validation
4. Add audit logging for sensitive operations

### Input Validation & Sanitization

**Request Validation**:
1. Add comprehensive validation rules
2. Implement custom validators
3. Sanitize user inputs
4. Validate file uploads

**SQL Injection Prevention**:
1. Use parameterized queries (sqlc)
2. Validate dynamic SQL if used
3. Implement query allowlists
4. Regular security audits

---

### 7. Deployment & DevOps Tasks

### Containerization

**Docker Optimization**:
1. Use multi-stage builds
2. Minimize image size
3. Implement proper health checks
4. Use non-root users

**Environment Configuration**:
1. Environment-specific configs
2. Secret management
3. Configuration validation
4. Graceful configuration reloads

### CI/CD Pipeline

**Quality Gates**:
1. All tests must pass
2. Code coverage requirements
3. Security vulnerability scanning
4. Linting and formatting checks

**Deployment Automation**:
1. Automated database migrations
2. Blue-green deployments
3. Rollback procedures
4. Health check validations

---

### 8. Monitoring & Observability Tasks

### Logging Enhancement

**Structured Logging**:
1. Add context to all log entries
2. Implement correlation IDs
3. Use consistent log levels
4. Add performance logging

**Error Tracking**:
1. Implement error categorization
2. Add stack traces for errors
3. Monitor error rates
4. Set up alerting for critical errors

### Metrics Collection

**Application Metrics**:
1. Request/response metrics
2. Database operation metrics
3. Cache hit/miss ratios
4. Business metrics tracking

**Infrastructure Metrics**:
1. Resource utilization monitoring
2. Database performance metrics
3. Network latency monitoring
4. Container health monitoring

---

## Task Templates

### Bug Fix Template

**Description**: Fix for [issue description]

**Steps**:
1. Reproduce the issue
2. Identify root cause
3. Write failing test
4. Implement fix
5. Verify test passes
6. Add regression tests
7. Update documentation if needed

**Verification**:
- [ ] Bug is fixed
- [ ] Tests pass
- [ ] No regressions
- [ ] Documentation updated

### Feature Development Template

**Description**: Implementation of [feature name]

**Requirements**:
- [ ] Domain model defined
- [ ] Database migration created
- [ ] Repository implemented
- [ ] Usecase implemented
- [ ] Handler implemented
- [ ] Routes registered
- [ ] Tests written
- [ ] Documentation updated

**Acceptance Criteria**:
- [ ] All requirements met
- [ ] Tests pass with >90% coverage
- [ ] Documentation is complete
- [ ] Code review approved

### Performance Improvement Template

**Description**: Performance optimization for [component]

**Baseline Metrics**:
- Response time: Xms
- Memory usage: YMB
- CPU usage: Z%
- Database queries: N

**Optimization Steps**:
1. Profile current performance
2. Identify bottlenecks
3. Implement optimizations
4. Measure improvements
5. Validate no regressions

**Success Criteria**:
- [ ] Response time improved by X%
- [ ] Memory usage reduced by Y%
- [ ] No functional regressions
- [ ] Tests pass

---

## Automation Scripts

### Pre-commit Hooks
```bash
#!/bin/sh
# .git/hooks/pre-commit

# Run tests
make test || exit 1

# Run linter
make lint || exit 1

# Check code generation
make sqlc
make generate-mocks

# Check if files changed
if git diff --quiet; then
    exit 0
else
    echo "Generated files changed. Please commit them."
    exit 1
fi
```

### Release Script
```bash
#!/bin/bash
# scripts/release.sh

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

# Run full CI
make ci

# Tag release
git tag -a "v$VERSION" -m "Release v$VERSION"
git push origin "v$VERSION"

# Build release binaries
make build

# Generate release notes
make generate-release-notes VERSION=$VERSION
```

---

## Troubleshooting Common Issues

### Database Migration Issues

**Problem**: Migration fails halfway through
**Solution**:
1. Check migration version: `make migrate-version`
2. Identify failed migration
3. Manually rollback if needed
4. Fix migration script
5. Re-run migration

### Test Failures

**Problem**: Tests fail with database connection errors
**Solution**:
1. Check test database configuration
2. Verify test database is running
3. Check migration status on test database
4. Verify test fixtures are correct

### Code Generation Issues

**Problem**: sqlc generation fails
**Solution**:
1. Check SQL syntax in query files
2. Verify migration files are valid
3. Check sqlc.yaml configuration
4. Ensure database schema is up to date

### Performance Issues

**Problem**: Slow API responses
**Solution**:
1. Check database query performance
2. Analyze application logs
3. Profile memory usage
4. Check connection pool settings

---

## Best Practices Checklist

### Code Quality
- [ ] All code follows Go conventions
- [ ] Functions have appropriate error handling
- [ ] Interfaces are used for dependencies
- [ ] Tests cover all public methods
- [ ] Documentation is up to date

### Security
- [ ] Input validation is implemented
- [ ] Authentication is properly configured
- [ ] Authorization checks are in place
- [ ] Sensitive data is not logged
- [ ] Dependencies are regularly updated

### Performance
- [ ] Database queries are optimized
- [ ] Connection pooling is configured
- [ ] Caching is used appropriately
- [ ] Memory usage is monitored
- [ ] Response times are acceptable

### Operations
- [ ] Health checks are implemented
- [ ] Logging is structured and informative
- [ ] Metrics are collected
- [ ] Alerts are configured
- [ ] Documentation is complete