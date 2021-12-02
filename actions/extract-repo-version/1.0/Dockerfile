FROM registry.erda.cloud/erda/terminus-alpine:base

RUN echo "http://mirrors.aliyun.com/alpine/v3.12/main/" > /etc/apk/repositories && \
    echo "http://mirrors.aliyun.com/alpine/v3.12/community/" >> /etc/apk/repositories && \
    apk update && apk add --no-cache git

COPY actions/extract-repo-version/1.0/run.sh /opt/action/run
RUN chmod +x /opt/action/run
