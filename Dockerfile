# Build stage
FROM mirror.gcr.io/library/golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api

# Install migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Final stage
FROM mirror.gcr.io/library/alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/db ./db

EXPOSE 8080

CMD ["/app/main"]
