# source: https://github.com/grafana/pyroscope-java/releases
# pyroscope.jar natively supports amd64/arm64 arch.
# you can use the following:
# - COPY --from=pyroscope-java /app/pyroscope-java.jar /your_target_path

ARG DOCKER_REGISTRY="registry.erda.cloud"

FROM ${DOCKER_REGISTRY}/retag/curl:8.0.1-alpine3.18.0 as curl
ARG PYROSCOPE_JAVA_VERSION="v0.11.3"

WORKDIR /app

RUN curl -OL https://github.com/grafana/pyroscope-java/releases/download/${PYROSCOPE_JAVA_VERSION}/pyroscope.jar

FROM ${DOCKER_REGISTRY}/retag/alpine:3.18.0

WORKDIR /app

COPY --from=curl /app .