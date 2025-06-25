FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/buildkit:v0.9.2 as buildkit
FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/golang:1.22-bookworm as builder

MAINTAINER shenli shenli@terminus.io

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions
RUN mkdir -p /opt/action/comp && \
    cp -r actions/php/1.0/comp/* /opt/action/comp

ARG ARCH

# go build
RUN GOOS=linux GOARCH=$ARCH go build -o /assets/run /go/src/github.com/erda-project/erda-actions/actions/php/1.0/internal/cmd/main.go

FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/php:cli-bookworm

COPY --from=composer /usr/bin/composer /usr/bin/composer
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl

# Installation of dependent tools
RUN apt-get update && apt-get install -y \
    unzip git make docker.io libzip-dev libpng-dev libjpeg-dev libonig-dev libxml2-dev libcurl4-openssl-dev libpq-dev \
    && rm -rf /var/lib/apt/lists/*
# Install PHP extensions
RUN docker-php-ext-install pdo pdo_mysql mbstring xml gd curl zip pgsql pdo_pgsql

# Install the Redis extension
RUN pecl install redis && docker-php-ext-enable redis

ENV COMPOSER_ALLOW_SUPERUSER=1

RUN composer config -g repo.packagist composer https://mirrors.tencent.com/composer/
RUN composer global config --no-plugins allow-plugins.slince/composer-registry-manager true
RUN composer global config --no-plugins allow-plugins.composer/installers true
RUN composer global require slince/composer-registry-manager
RUN composer repo:use tencent
RUN composer config -g process-timeout 600
COPY --from=builder /assets /opt/action
COPY --from=builder /opt/action/comp /opt/action/comp
