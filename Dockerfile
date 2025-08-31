# Build stage
FROM golang:alpine AS builder

# Install git and ca-certificates (needed for go mod download and HTTPS requests)
RUN apk add --no-cache git ca-certificates tzdata tini-static

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o server \
    cmd/server/main.go

# Final stage
FROM alpine

# Copy tini from builder stage
COPY --from=builder /sbin/tini-static /sbin/tini-static

# Copy ca-certificates from builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data from builder stage
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary from builder stage
COPY --from=builder /build/server /server

# Copy configuration files
COPY --from=builder /build/configs /configs

EXPOSE 8080

ENTRYPOINT [ "tini-static --" ]
CMD ["/server"]