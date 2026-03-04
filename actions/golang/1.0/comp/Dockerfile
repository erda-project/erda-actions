FROM registry.erda.cloud/erda-x/debian-bookworm:12

ARG TARGET

COPY assets /opt/service/
COPY ${TARGET} /opt/service/run

RUN chmod +x /opt/service/run
WORKDIR /opt/service/
CMD ["/opt/service/run"]