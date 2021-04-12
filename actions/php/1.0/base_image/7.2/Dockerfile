FROM php:7.2-apache
RUN pecl install redis-5.1.1 \
    && docker-php-ext-enable redis
RUN a2enmod rewrite