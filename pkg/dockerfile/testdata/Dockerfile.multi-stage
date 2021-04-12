ARG VERSION
FROM alpine:${VERSION} AS builder
RUN echo hello
ARG DEP_CMD
RUN eval ${DEP_CMD}

FROM alpine:3.7
COPY --from=builder /etc/hosts /tmp/hosts
RUN cat /tmp/hosts

FROM centos:7