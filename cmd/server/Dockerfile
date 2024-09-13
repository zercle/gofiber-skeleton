# syntax=docker/dockerfile:1
FROM golang AS builder
ARG GITHUB_USER_BUILD
ARG GITHUB_TOKEN_BUILD
ARG timezone=Asia/Bangkok
ENV GITHUB_USER $GITHUB_USER_BUILD
ENV GITHUB_TOKEN $GITHUB_TOKEN_BUILD
ENV TZ $timezone
WORKDIR /app

# install tini
RUN apt-get update && \
    apt-get -y install tini && \
    apt-get -y clean

# for CI/CD with private repo
# RUN go env -w GOPRIVATE=github.com/zercle

# use ssh instead of https
# RUN git config --global url."git@github.com:".insteadOf "https://github.com"

# use https with github token
# RUN git config --global url."https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com".insteadOf "https://github.com"

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.* ./
RUN go mod download && go mod verify

# let's build project
COPY . .
RUN go test -v ./... && \
    CGO_ENABLED=0 go build -v \
    -installsuffix 'static' \
    -ldflags="-X 'main.version=$(git rev-parse --short HEAD)' -X 'main.build=$(date --iso-8601=seconds)'" \
    -o dist/server ./cmd/server

# pack PRD image
FROM gcr.io/distroless/base:nonroot
LABEL maintainer="Kawin Viriyaprasopsook <kawin.vir@zercle.tech>"

ARG timezone=Asia/Bangkok

ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV TZ $timezone

# Create app dir
WORKDIR /app
COPY --from=builder /usr/bin/tini-static /usr/bin/tini
COPY --from=builder /app/dist/server /app/server
COPY ./configs /app/configs
COPY ./internal/assets /app/internal/assets
COPY ./web /app/web

EXPOSE 8080 8443

# default run entrypoint
ENTRYPOINT ["tini", "--", "/app/server"]
CMD ["--env=prd"]
