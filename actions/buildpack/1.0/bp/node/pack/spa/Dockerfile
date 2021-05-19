FROM {{BP_DOCKER_BASE_REGISTRY}}/erda/dice-nginx:1.1.0

# Set special timezone
RUN echo "Asia/Shanghai" | tee /etc/timezone

ARG PUBLIC_DIR="public/"

ADD /app /app

COPY /app/${PUBLIC_DIR} /usr/share/nginx/html
# TODO add default template
RUN cp /app/nginx.conf.template /etc/nginx/conf.d/nginx.conf.template

CMD sed -i "s^server_name .*^^g" /etc/nginx/conf.d/nginx.conf.template && \
    envsubst "`printf '$%s' $(bash -c "compgen -e")`" < /etc/nginx/conf.d/nginx.conf.template > /etc/nginx/conf.d/default.conf && \
    nginx -g 'daemon off;'
