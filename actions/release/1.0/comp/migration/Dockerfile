# Build app image
FROM registry.cn-hangzhou.aliyuncs.com/dice-third-party/migration:1.0.2

RUN echo "Asia/Shanghai" | tee /etc/timezone

COPY sql/ /tmp/db/migration
RUN chmod +x /entrypoint.sh

WORKDIR /

CMD ["/entrypoint.sh"]