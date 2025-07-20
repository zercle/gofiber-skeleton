# syntax=docker/dockerfile:1

# Build stage
FROM mirror.gcr.io/library/golang AS build

WORKDIR /app

COPY go.mod go.sum ./ 
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o dist/server ./cmd/api

# Tini stage
FROM mirror.gcr.io/library/debian:stable-slim AS pkg
RUN apt-get update && apt-get install -y tini curl && cp $(which curl) /usr/bin/curl-static

# Final stage
FROM gcr.io/distroless/base:nonroot

WORKDIR /app

COPY --from=build /app/dist/server server
COPY --from=pkg /usr/bin/tini-static /usr/bin/tini

EXPOSE 8080

ENTRYPOINT ["/usr/bin/tini", "--", "/app/server"]
