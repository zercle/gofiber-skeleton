# GitHub Actions Workflows

This directory contains all CI/CD workflows for the project. The workflows are organized to minimize duplication and provide fast feedback while ensuring comprehensive testing.

## Workflow Overview

### 1. CI (`ci.yml`)
**Purpose:** Fast feedback on every push and pull request

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main`

**Jobs:**
- **Lint**: Runs golangci-lint v2 for code quality checks
- **Quick Tests**: Unit tests with PostgreSQL service (`-short` flag)
- **Security Scan**: Gosec and Trivy security scanning

**Run Time:** ~3-5 minutes

**Use Case:** Primary workflow for all PRs and commits. Provides quick feedback to developers.

---

### 2. Integration Tests (`integration.yml`)
**Purpose:** Comprehensive testing with full infrastructure

**Triggers:**
- Push to `main` branch
- Manual trigger (`workflow_dispatch`)

**Jobs:**
- **Full Integration Tests**: Complete test suite with:
  - Code generation (sqlc, swag, mockgen)
  - Docker Compose services (PostgreSQL, Valkey)
  - Database migrations
  - Full integration tests
  - Application build

**Run Time:** ~8-12 minutes

**Use Case:** Validates that all components work together correctly. Runs after merging to main or can be triggered manually for testing.

---

### 3. CD (`cd.yml`)
**Purpose:** Build and deploy application

**Triggers:**
- Git tags matching `v*` pattern (e.g., `v1.0.0`)

**Jobs:**
- **Build and Push**: Docker image build and push to GitHub Container Registry
- **Deploy Staging**: Automatic deployment to staging (on main branch push)
- **Deploy Production**: Manual deployment to production (on version tags)
- **Notify**: Deployment status notifications

**Run Time:** ~5-10 minutes (+ deployment time)

**Use Case:** Automated releases and deployments when version tags are created.

---

### 4. Setup Workflow (`_setup.yml`)
**Purpose:** Reusable workflow for common Go setup tasks

**Status:** Currently not used, available for future workflow composition

**Features:**
- Go environment setup with caching
- Optional tool installation (sqlc, swag, mockgen)
- Optional code generation

---

## Workflow Strategy

### Fast Feedback Loop
```
Developer Push → CI Workflow (3-5 min)
├─ Lint (parallel)
├─ Quick Tests (parallel)
└─ Security Scan (parallel)
```

### Comprehensive Validation
```
Merge to Main → Integration Workflow (8-12 min)
└─ Full integration tests with real services
```

### Release Pipeline
```
Create Tag → CD Workflow (5-10 min)
├─ Build Docker Image
├─ Deploy to Staging (if main branch)
└─ Deploy to Production (manual approval)
```

---

## Configuration Details

### Environment Variables

**Integration Tests:**
```yaml
DB_HOST: localhost
DB_PORT: 5432
DB_USER: postgres
DB_PASSWORD: postgres
DB_NAME: gofiber_skeleton
VALKEY_HOST: localhost
VALKEY_PORT: 6379
JWT_SECRET: test-secret-key-for-integration-tests
```

### Golangci-lint Configuration
- Version: v2 (latest major version)
- Timeout: 5 minutes
- Configuration: `.golangci.yml` in project root

### Code Generation Tools
- `sqlc`: SQL to Go code generation
- `swag`: Swagger documentation generation
- `mockgen`: Mock generation for testing

---

## Concurrency Control

All workflows use concurrency groups to prevent multiple runs:
```yaml
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true
```

This ensures:
- Only one workflow runs per branch at a time
- New pushes cancel old, in-progress runs
- Saves CI/CD resources

---

## Best Practices

### For Developers

1. **Before committing:**
   - Run `golangci-lint run` locally
   - Run `go test ./...` to catch basic issues

2. **During PR review:**
   - CI workflow must pass (all 3 jobs green)
   - Address all lint warnings
   - Ensure security scan passes

3. **Before merging:**
   - Ensure PR is up-to-date with main
   - All CI checks passing
   - Code reviewed and approved

### For Releases

1. **Creating a release:**
   ```bash
   git tag -a v1.2.3 -m "Release version 1.2.3"
   git push origin v1.2.3
   ```

2. **Manual integration test:**
   - Go to Actions → Integration Tests → Run workflow

3. **Monitoring deployments:**
   - Check Actions tab for CD workflow status
   - Verify staging deployment before production

---

## Troubleshooting

### CI Workflow Fails

**Lint failures:**
- Run `golangci-lint run` locally to see full output
- Fix issues or add exclusions to `.golangci.yml` if needed

**Test failures:**
- Check test logs in GitHub Actions
- Tests use PostgreSQL service, ensure migrations are compatible

**Security scan failures:**
- Review Gosec and Trivy findings
- Address vulnerabilities or add exclusions if false positives

### Integration Workflow Fails

**Code generation fails:**
- Ensure `sqlc.yaml` and swagger comments are correct
- Run `sqlc generate` and `swag init` locally

**Services not ready:**
- Check Docker Compose configuration
- Ensure health checks are properly configured

**Migration failures:**
- Verify migration files are correct
- Check database connection configuration

---

## Maintenance

### Updating Workflows

1. **Add new linters:**
   - Update `.golangci.yml`
   - Test locally with `golangci-lint run`

2. **Add new services:**
   - Update `docker-compose.yml`
   - Add service startup/health checks in `integration.yml`

3. **Change Go version:**
   - Update `go-version: 'stable'` in all workflows
   - Or pin to specific version: `go-version: '1.23'`

### Monitoring Performance

- Check workflow run times in Actions tab
- Optimize slow jobs by:
  - Using Go build cache effectively
  - Running jobs in parallel where possible
  - Reducing test verbosity in quick tests

---

## Security Considerations

- Secrets are managed via GitHub Secrets
- Never commit credentials or secrets to workflows
- Use least-privilege permissions for each job
- Container images scanned with Trivy before deployment
- Dependencies regularly updated via Dependabot
