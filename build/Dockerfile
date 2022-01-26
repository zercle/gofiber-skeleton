# syntax=docker/dockerfile:1
FROM golang AS builder
ARG GITHUB_USER_BUILD
ARG GITHUB_TOKEN_BUILD
ENV GITHUB_USER $GITHUB_USER_BUILD
ENV GITHUB_TOKEN $GITHUB_TOKEN_BUILD
WORKDIR /src/app
RUN go env -w GOPRIVATE=github.com/zercle
RUN go get -d -v -u
# RUN git config --global url."git@github.com:".insteadOf "https://github.com"
# RUN git config --global url."https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com".insteadOf "https://github.com"
COPY . .
RUN go mod tidy && go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'main.version=$(git rev-parse --short HEAD)' -X 'main.build=$(date --iso-8601=seconds)Z'" -a -installsuffix cgo -mod vendor -o app-dist .

FROM alpine
LABEL maintainer="Kawin Viriyaprasopsook <kawin.vir@zercle.tech>"

ARG	timezone=Asia/Bangkok

ENV	LANG en_US.UTF-8
ENV	LC_ALL en_US.UTF-8
ENV	TZ $timezone

# Add config repositories
RUN	echo 'https://download.nus.edu.sg/mirror/alpine/latest-stable/main' > /etc/apk/repositories \
    && echo 'https://download.nus.edu.sg/mirror/alpine/latest-stable/community' >> /etc/apk/repositories \
    # && echo '@edge https://download.nus.edu.sg/mirror/alpine/edge/main' >> /etc/apk/repositories \
    # && echo '@testing https://download.nus.edu.sg/mirror/alpine/edge/testing' >> /etc/apk/repositories \
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
COPY --from=builder /src/app/app-dist /app/
COPY *.env /app/
COPY build/docker/ /app/
COPY assets /app/assets
COPY internal/assets /app/internal/assets
RUN ln -sf /app/app-dist /usr/local/bin/app \
    && chmod +x /app/docker-entrypoint.sh
WORKDIR /app

ENTRYPOINT ["/app/docker-entrypoint.sh"]

CMD ["app", "env=dev"]

EXPOSE 8080 8443
