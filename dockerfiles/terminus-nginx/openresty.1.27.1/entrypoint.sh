#!/bin/sh
set -e
if [[ -f /etc/nginx/conf.d/nginx.conf.override ]]; then
    rm /usr/local/openresty/nginx/conf/nginx.conf
    cp /etc/nginx/conf.d/nginx.conf.override /usr/local/openresty/nginx/conf/nginx.conf
fi
exec "$@"
