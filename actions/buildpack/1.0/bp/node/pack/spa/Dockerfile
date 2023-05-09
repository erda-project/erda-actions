FROM --platform=$TARGETPLATFORM {{BP_DOCKER_BASE_REGISTRY}}/erda-addons/openresty:1.21.4

# Set special timezone
RUN echo "Asia/Shanghai" | tee /etc/timezone

ARG PUBLIC_DIR="public"

ADD /app /app

RUN mkdir -p /usr/share/nginx/html && cp -r /app/${PUBLIC_DIR}/* /usr/share/nginx/html
# TODO add default template
RUN cp /app/nginx.conf.template /etc/nginx/conf.d/nginx.conf.template

CMD sed -i "s^server_name .*^^g" /etc/nginx/conf.d/nginx.conf.template && \
    envsubst "`printf '$%s' $(bash -c "compgen -e")`" < /etc/nginx/conf.d/nginx.conf.template > /etc/nginx/conf.d/default.conf && \
    nginx -g 'daemon off;'
