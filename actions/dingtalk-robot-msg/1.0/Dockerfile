FROM registry.erda.cloud/erda/alpine:3.14

RUN echo "http://mirrors.aliyun.com/alpine/v3.14/main/" > /etc/apk/repositories && \
    echo "http://mirrors.aliyun.com/alpine/v3.14/community/" >> /etc/apk/repositories && \
    apk update && apk add --no-cache openssl curl jsonnet

COPY actions/dingtalk-robot-msg/1.0/run.sh /opt/action/run
RUN chmod +x /opt/action/run
