# syntax=docker/dockerfile:1

FROM golang AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o server ./cmd/server/main.go

# ---

FROM gcr.io/distroless/base AS runner

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY config/config.yaml /app/config.yaml
COPY migrations /app/migrations

EXPOSE 8080

ENV PORT=8080

ENTRYPOINT ["/app/server"]
