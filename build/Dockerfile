# syntax=docker/dockerfile:1
FROM golang AS builder
ARG GITHUB_USER_BUILD
ARG GITHUB_TOKEN_BUILD
ENV GITHUB_USER $GITHUB_USER_BUILD
ENV GITHUB_TOKEN $GITHUB_TOKEN_BUILD
WORKDIR /usr/src/app

# for CI/CD with private repo
# RUN go env -w GOPRIVATE=github.com/zercle
# RUN git config --global url."git@github.com:".insteadOf "https://github.com"
# RUN git config --global url."https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com".insteadOf "https://github.com"

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# let's build project
COPY . .
RUN go build -v -buildvcs=false -ldflags="-X 'main.version=$(git rev-parse --short HEAD)' -X 'main.build=$(date --iso-8601=seconds)Z'" -o app-dist .

FROM alpine
LABEL maintainer="Kawin Viriyaprasopsook <kawin.vir@zercle.tech>"

ARG	timezone=Asia/Bangkok

ENV	LANG en_US.UTF-8
ENV	LC_ALL en_US.UTF-8
ENV	TZ $timezone

# Add config repositories
RUN	echo 'https://dl-cdn.alpinelinux.org/alpine/latest-stable/main' > /etc/apk/repositories \
    && echo 'https://dl-cdn.alpinelinux.org/alpine/latest-stable/community' >> /etc/apk/repositories \
    && mkdir /run/openrc \
    && touch /run/openrc/softlevel

# Update OS
RUN	apk update && apk upgrade \
    && apk add --no-cache \
    openrc \
    tzdata \
    bash \
    bash-completion \
    net-tools \
    bind-tools \
    mtr \
    ca-certificates

# Change locale
RUN echo $timezone > /etc/timezone \
    && cp /usr/share/zoneinfo/$timezone /etc/localtime

# Create app dir
RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /usr/src/app/app-dist /app/
COPY *.env /app/
COPY build/docker/ /app/
COPY assets /app/assets
COPY internal/assets /app/internal/assets
RUN ln -sf /app/app-dist /usr/local/bin/app && chmod +x /app/docker-entrypoint.sh

EXPOSE 8080 8443

# startup script
ENTRYPOINT ["bash", "/app/docker-entrypoint.sh"]
# fallback
CMD ["app", "env=dev"]
