# Build app image
FROM erda-registry.cn-hangzhou.cr.aliyuncs.com/erda/migration:1.0.2

RUN echo "Asia/Shanghai" | tee /etc/timezone

COPY sql/ /tmp/db/migration
RUN chmod +x /entrypoint.sh

WORKDIR /

CMD ["/entrypoint.sh"]