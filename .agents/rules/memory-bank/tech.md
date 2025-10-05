# Technology Stack

This document enumerates the sanctioned tooling, libraries, and operational practices for the modern Go Fiber template. All components are selected for long-term maintainability, active community support, and compatibility with cloud-native deployments.

## Languages and Runtime

- **Go 1.25+** with toolchain pinning via `go env -w GOTOOLCHAIN=go1.25.0`.
- **TypeScript 5.x** for supporting scripts, infrastructure CDK, and optional frontend SDKs.
- **SQL** (PostgreSQL dialect) with sqlc code generation.

## Core Frameworks & Libraries

| Concern | Technology |
| --- | --- |
| HTTP server | Fiber v3 with middleware ecosystem (recover, cors, limiter, helmet) |
| gRPC services | `google.golang.org/grpc` with Buf toolchain |
| Configuration | Viper + env overrides, SOPS-encrypted config bundles |
| Dependency injection | Uber Fx or Wire for compile-time graph |
| Validation | go-playground/validator v10 with domain-specific wrappers |
| Auth & Security | PASETO v2 tokens, bcrypt/argon2id hashing, OPA for policy checks |
| Background jobs | Asynq (Redis-backed) for scheduled workflows |
| Rate limiting | Redis-based sliding window implemented via Go Redis client |

## Data Layer

- **Primary database**: PostgreSQL 18 (Cloud SQL/Aurora compatible) using **sqlc** generated repositories, migrations managed by **Atlas**.
- **Caching**: Valkey 8 via go-redis client.
- **Search**: Meilisearch 1.x with HTTP SDK integration.
- **Messaging/Eventing**: NATS JetStream for pub/sub and command bus semantics.
- **Analytics pipeline**: ClickHouse via Kafka Connect (optional module).

## Observability

- **Tracing**: OpenTelemetry Go SDK exporting OTLP → Collector → Tempo/Jaeger.
- **Metrics**: Prometheus client with exemplars, Grafana dashboards provisioned via Jsonnet.
- **Logging**: Zerolog structured logs with trace correlation; Loki integration through Promtail sidecars.
- **Error tracking**: Sentry Go SDK with performance tracing enabled.

## DevEx Tooling

- **Package management**: Go modules (proxy-aware), npm/pnpm for TypeScript assets.
- **Task runner**: Taskfile.yml orchestrating dev, build, test, and codegen flows.
- **Hot reload**: Air for Go binaries, Tilt for multi-service dev orchestration.
- **Linting & formatting**: golangci-lint, gofmt, revive, sqlfluff (SQL), eslint/prettier (TS).
- **Code generation**: sqlc, Buf (proto), mockery (interfaces), oapi-codegen (OpenAPI), swagger docs via Swag.
- **Secrets management**: SOPS + age, integrated with AWS/GCP KMS.

## Testing Strategy

- Unit tests with Go testing package + Testify.
- Mock generation via mockery and go.uber.org/mock.
- Integration tests leveraging Testcontainers-Go for ephemeral PostgreSQL/Redis/NATS.
- Contract tests generated from OpenAPI/gRPC definitions validated with Dredd/Buf.
- Load testing harness using k6 scripts stored under `tools/k6`.

## CI/CD

- GitHub Actions workflow matrix:
  - Lint & unit tests (Go, SQL, Proto).
  - Integration tests using service containers.
  - Build & push multi-arch Docker images (linux/amd64, linux/arm64).
  - Upload coverage to Codecov.
  - Trigger Helm chart packaging via `helmfile sync`.
- Deployment promotion orchestrated by Argo CD with progressive delivery (Argo Rollouts).
- Security scans: Trivy (images), Gosec (code), Dependabot updates.

## Infrastructure as Code

- **Terraform** modules provisioning cloud infrastructure (networking, managed DB/cache, secrets).
- **Helmfile** manages Kubernetes releases; base charts align with official community charts.
- **Cilium** for service mesh & network security, optional Linkerd integration for lightweight scenarios.
- **External Secrets Operator** to bridge Vault/Secrets Manager into cluster secrets.

## Compliance & Governance

- Policy enforcement via Open Policy Agent (OPA) integrated in API gateway and CI checks.
- Data retention & residency handled through PostgreSQL partitioning strategies and policy configs.
- Audit trails stored in immutable PostgreSQL schemas with time-based access controls.
- Static application security testing (SAST) and dependency CVE scanning part of CI pipeline.

## Supported Environments

| Environment | Runtime | Tooling |
| --- | --- | --- |
| Local | Docker Desktop/Colima with Tilt, Air, Taskfile | Compose, Meilisearch, NATS, Redis, Postgres |
| Staging | Kubernetes (K3s/KinD or managed) | Helmfile, Argo CD, Prometheus stack |
| Production | Managed Kubernetes (GKE/AKS/EKS) | Cloud SQL/Aurora, Managed Redis/NATS, Cloud Load Balancers, Cloud Armor/WAF |

## Reference Repositories & Inspirations

- gofiber-boilerplate community templates (architecture inspiration)
- Entgo clean architecture example for domain layering patterns
- OpenTelemetry demo for observability wiring best practices
- Buffer/Segment public infra blueprints for data streaming and governance