FROM {{CENTRAL_REGISTRY}}/erda/terminus-php-apache:{{PHP_VERSION}}
ARG TARGET

RUN sed -ri -e 's!/var/www/html!{{APACHE_DOCUMENT_ROOT}}!g' /etc/apache2/sites-available/*.conf
RUN sed -ri -e 's!/var/www/!{{APACHE_DOCUMENT_ROOT}}!g' /etc/apache2/apache2.conf /etc/apache2/conf-available/*.conf
RUN a2enmod rewrite

COPY ${TARGET} /var/www/html/

CMD ["apache2-foreground"]