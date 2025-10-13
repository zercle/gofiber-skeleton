# Build stage
FROM golang:alpine AS builder

ENV TZ=Asia/Bangkok
ENV LANG=C.UTF-8

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata build-base

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server cmd/server/main.go

# Final stage
FROM alpine:latest

ENV TZ=Asia/Bangkok
ENV LANG=C.UTF-8

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata tini

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Copy configuration files
COPY --from=builder /app/config.yaml .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Execute the init command
ENTRYPOINT [ "/sbin/tini", "--" ]

# Run the application
CMD [ "./server" ]