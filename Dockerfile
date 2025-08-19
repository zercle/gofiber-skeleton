FROM golang:1.24.6-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/air-verse/air@latest

COPY . .

RUN go build -o main ./cmd/server

EXPOSE 8080

CMD ["./main"]