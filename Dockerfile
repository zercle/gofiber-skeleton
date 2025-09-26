# Build Stage
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w -extldflags '-static'" -o /app/bin/server ./cmd/server

# Run Stage
FROM alpine

WORKDIR /app

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

COPY --from=builder /app/bin/server /app/bin/server
COPY --from=builder /app/.env.example /app/.env.example

# Expose the port the app runs on
EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD ["/app/bin/server", "health"]

# Command to run the executable
CMD ["/app/bin/server"]