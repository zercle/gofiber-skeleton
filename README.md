# GoFiber Booking Store

A boilerplate Go project using Clean Architecture and SOLID principles for a booking store API. Built with [Fiber](https://github.com/gofiber/fiber), [GORM](https://gorm.io/) (PostgreSQL), JWT authentication, Viper configuration, and graceful shutdown support.

## Features

- Clean Architecture: separation of domain, usecases, repositories, handlers, and middleware
- SOLID design: dependency inversion via interfaces, single-responsibility components
- Fiber HTTP server with graceful shutdown
- JWT-based authentication middleware
- PostgreSQL integration using GORM
- Configuration via Viper with environment variable overrides
- Modular, testable components
- Database migrations directory for schema setup
- GitHub Actions for CI (testing, linting)

## Tech Stack

- Go 1.24+
- Fiber v2
- GORM v2 with PostgreSQL driver
- JWT (github.com/golang-jwt/jwt/v5)
- Viper for config
- GitHub Actions (go.yml, golangci-lint.yml)

## Getting Started

### Prerequisites

- Go 1.24 or newer
- PostgreSQL database
- Git


### Configuration



Create a `config.yaml` in the project root or set environment variables:



```yaml

PORT: "8080"

DATABASE_URL: "postgres://user:password@localhost:5432/bookingdb?sslmode=disable"

JWT_SECRET: "your-secret-key"

```



Alternatively export environment variables:



```

export PORT=8080

export DATABASE_URL="postgres://user:password@localhost:5432/bookingdb?sslmode=disable"

export JWT_SECRET="your-secret-key"

```


### Directory Structure

```
gofiber-skeleton/
├── cmd/server/           # main.go entrypoint
├── config/               # Viper config loader
├── internal/
│   ├── domain/           # Entities & repository interfaces
│   ├── usecase/          # Business logic
│   ├── repository/       # GORM implementations
│   ├── handler/          # HTTP route handlers
│   └── middleware/       # JWT auth, error handlers
├── pkg/jwtutil/          # JWT utility functions
├── migrations/           # SQL migration files
├── go.mod
└── README.md
```


## Installation & Running



```bash

git clone https://github.com/your-org/gofiber-skeleton.git

cd gofiber-skeleton



# install dependencies

go mod download



# run migrations (example using plain psql)

psql $DATABASE_URL -f migrations/001_create_bookings.sql



# start server

go run cmd/server/main.go

```



Server will listen on `http://localhost:${PORT}`.


## JWT Ed25519 Key Generation

This project uses Ed25519 keys for JWT signing and verification.

### Generate Ed25519 Key Pair

You can generate a key pair using OpenSSL or Go:

**With OpenSSL:**
```bash
openssl genpkey -algorithm Ed25519 -out ed25519-private.pem
openssl pkey -in ed25519-private.pem -pubout -out ed25519-public.pem
```

**Convert to Base64 (for config):**
```bash
# Private key (base64, no headers/footers)
awk 'NF {sub(/-----.*-----/, ""); print}' ed25519-private.pem | base64 -w0
# Public key (base64, no headers/footers)
awk 'NF {sub(/-----.*-----/, ""); print}' ed25519-public.pem | base64 -w0
```

**With Go:**
```go
package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

func main() {
	pub, priv, _ := ed25519.GenerateKey(nil)
	fmt.Println("Private:", base64.StdEncoding.EncodeToString(priv))
	fmt.Println("Public: ", base64.StdEncoding.EncodeToString(pub))
}
```

Set these base64 strings as `JWT_PRIVATE_KEY` and `JWT_PUBLIC_KEY` in your config or environment.

## Docker & Docker Compose



You can run the API and a PostgreSQL database using Docker Compose.


### 1. Build and Run with Docker Compose

```bash
docker-compose up --build
```

This will start both the API and a PostgreSQL database. The API will be available at [http://localhost:8080](http://localhost:8080).

### 2. Environment Variables

You can override environment variables in the `docker-compose.yml` file or by creating a `.env` file.

### 3. Running Migrations

You can run the migration manually (from your host or inside the container):

```bash
docker-compose exec db psql -U postgres -d bookingdb -f /migrations/001_create_bookings.sql
```

Or use a migration tool as needed.



## Graceful Shutdown

The server listens for `SIGINT`/`SIGTERM` and shuts down with a 5-second timeout for active requests to complete.

## Testing

```bash
go test ./... -v -cover
```


## CI / Linting



Configured workflows under `.github/workflows`:



- `go.yml`: runs tests on push and PR

- `golangci-lint.yml`: runs `golangci-lint` for code quality



## License


MIT © Your Name or Company  


MIT © Zercle Technology Co., Ltd.
