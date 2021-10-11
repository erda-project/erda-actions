# registry.erda.cloud/erda/terminus-openjdk:v1.8.0.242-filebeat.v6.7.0
FROM registry.erda.cloud/erda/terminus-openjdk:v1.8.0.242

ARG FILEBEAT_VERSION=6.7.0-linux-x86_64
ENV FILEBEAT_VERSION ${FILEBEAT_VERSION}
ARG FILEBEAT_DOWNLOAD_URL=https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-${FILEBEAT_VERSION}.tar.gz
ENV FILEBEAT_DOWNLOAD_URL ${FILEBEAT_DOWNLOAD_URL}

RUN mkdir -p /opt/filebeat
RUN curl -s ${FILEBEAT_DOWNLOAD_URL} | tar zx -C /tmp && \
    mv /tmp/filebeat-${FILEBEAT_VERSION}/* /opt/filebeat/
