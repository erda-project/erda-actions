FROM {{CENTRAL_REGISTRY}}/erda/terminus-php-apache:{{PHP_VERSION}}
ARG TARGET

# Supported version:
# - php 7.2 (base image: debian buster)
RUN sed -i -e 's|security.debian.org|mirrors.aliyun.com/debian-archive|g' \
           -e 's|deb.debian.org|mirrors.aliyun.com/debian-archive|g' /etc/apt/sources.list

RUN sed -ri -e 's!/var/www/html!{{APACHE_DOCUMENT_ROOT}}!g' /etc/apache2/sites-available/*.conf
RUN sed -ri -e 's!/var/www/!{{APACHE_DOCUMENT_ROOT}}!g' /etc/apache2/apache2.conf /etc/apache2/conf-available/*.conf
RUN a2enmod rewrite

RUN apt-get update && apt-get install -y \
    unzip git make docker.io libzip-dev libpng-dev libjpeg-dev libonig-dev libxml2-dev libcurl4-openssl-dev libpq-dev \
    && rm -rf /var/lib/apt/lists/*

RUN docker-php-ext-install pdo pdo_mysql mbstring xml gd curl zip pgsql pdo_pgsql

COPY ${TARGET} /var/www/html/

CMD ["apache2-foreground"]