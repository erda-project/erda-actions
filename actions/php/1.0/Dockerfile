FROM registry.erda.cloud/retag/buildkit:v0.11.3 AS buildkit
FROM registry.erda.cloud/retag/docker:20.10.24-cli AS docker-cli
FROM registry.erda.cloud/retag/composer:2.2.25 AS composer
FROM registry.erda.cloud/erda-x/golang:1.24 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions
RUN mkdir -p /opt/action/comp && \
    cp -r actions/php/1.0/comp/* /opt/action/comp

# go build
RUN go build -o /assets/run /go/src/github.com/erda-project/erda-actions/actions/php/1.0/internal/cmd/main.go

FROM registry.erda.cloud/retag/php:8.4.12-cli-bookworm

COPY --from=composer /usr/bin/composer /usr/bin/composer

# Installation of dependent tools
RUN sed -i 's|deb.debian.org|mirrors.aliyun.com|g' /etc/apt/sources.list.d/debian.sources && \
    apt-get update && apt-get install -y \
    unzip git make libzip-dev libpng-dev libjpeg-dev libonig-dev libxml2-dev libcurl4-openssl-dev libpq-dev \
    && rm -rf /var/lib/apt/lists/* \

# Install extensions
RUN docker-php-ext-install pdo pdo_mysql mbstring xml gd curl zip pgsql pdo_pgsql && \
    pecl install redis && docker-php-ext-enable redis

ENV COMPOSER_ALLOW_SUPERUSER=1

RUN composer config -g repo.packagist composer https://mirrors.tencent.com/composer/ && \
    composer global config --no-plugins allow-plugins.slince/composer-registry-manager true && \
    composer global config --no-plugins allow-plugins.composer/installers true && \
    composer global require slince/composer-registry-manager && \
    composer repo:use tencent && \
    composer config -g process-timeout 600

COPY --from=docker-cli /usr/local/bin/docker /usr/bin/docker
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
COPY --from=builder /assets /opt/action
COPY --from=builder /opt/action/comp /opt/action/comp
