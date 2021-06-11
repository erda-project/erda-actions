FROM registry.erda.cloud/erda/terminus-golang:1.14

ARG TARGET

RUN ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone

COPY assets /opt/service/
COPY ${TARGET} /opt/service/run

RUN chmod +x /opt/service/run
WORKDIR /opt/service/
CMD ["/opt/service/run"]