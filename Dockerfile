# syntax=docker/dockerfile:1
FROM golang AS builder

# Define build arguments at the top for clarity and reusability
ARG GITHUB_USER_BUILD
ARG GITHUB_TOKEN_BUILD
ARG timezone=Asia/Bangkok
ARG TINI_VERSION=0.19.0 # Specify tini version for consistent builds

# Set environment variables for the builder stage
ENV GITHUB_USER=${GITHUB_USER_BUILD}
ENV GITHUB_TOKEN=${GITHUB_TOKEN_BUILD}
ENV TZ=${timezone}

WORKDIR /app

# Download and install tini statically for the distroless image.
# This avoids apt-get in the builder stage and ensures compatibility with distroless.
RUN wget -qO /usr/local/bin/tini https://github.com/krallin/tini/releases/download/v${TINI_VERSION}/tini-amd64 && \
    chmod +x /usr/local/bin/tini

# for CI/CD with private repo
# RUN go env -w GOPRIVATE=github.com/zercle

# use ssh instead of https
# RUN git config --global url."git@github.com:".insteadOf "https://github.com"

# use https with github token
# RUN git config --global url."https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com".insteadOf "https://github.com"

# Pre-copy/cache go.mod and go.sum for efficient dependency management.
# This layer is only invalidated if go.mod or go.sum changes.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the application source code
COPY . .

# Build the Go application.
# - `go test -v ./...` is kept as per original, though often run in a separate CI stage.
# - `CGO_ENABLED=0` ensures a statically linked binary, suitable for distroless.
# - `-installsuffix 'static'` is redundant with `CGO_ENABLED=0` and removed.
# - `-ldflags` embeds version and build information into the binary.
# - Output binary is named `server` in `dist/server`.
RUN go test -v ./... && \
    CGO_ENABLED=0 go build -v \
    -ldflags="-X 'main.version=$(git rev-parse --short HEAD)' -X 'main.build=$(date --iso-8601=seconds)'" \
    -o dist/server ./cmd/app

# Pack PRD image
FROM gcr.io/distroless/base:nonroot
LABEL maintainer="Kawin Viriyaprasopsook <kawin.vir@zercle.tech>"

# Define timezone ARG again for the final stage, as ARGs are stage-scoped
ARG timezone=Asia/Bangkok

# Set environment variables for the final stage
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8
ENV TZ=${timezone}

# Create application directory
WORKDIR /app

# Copy the built application binary and the static tini binary from the builder stage.
# Corrected binary path from /app/main to /app/dist/server.
COPY --from=builder /app/dist/server /app/server
COPY --from=builder /usr/local/bin/tini /usr/bin/tini

# Copy configuration and database migrations
COPY configs ./configs
COPY database/migrations ./database/migrations

VOLUME /data

EXPOSE 3000

# Default run entrypoint using the copied static tini binary
ENTRYPOINT ["/usr/bin/tini", "--", "/app/server"]
