FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/terminus-golang:1.19.7

ARG TARGET

COPY assets /opt/service/
COPY ${TARGET} /opt/service/run

RUN chmod +x /opt/service/run
WORKDIR /opt/service/
CMD ["/opt/service/run"]