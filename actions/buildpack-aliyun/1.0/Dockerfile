FROM registry.erda.cloud/retag/buildkit:v0.11.3 AS buildkit
FROM registry.erda.cloud/retag/docker:28.2.2-cli AS docker-cli
FROM registry.erda.cloud/erda-x/golang:1.22 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

# go build
RUN go build -o /assets/run github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal

# async assets
WORKDIR actions/buildpack-aliyun/1.0
RUN bash sync_assets.sh

FROM registry.erda.cloud/erda-x/debian-bookworm:12

RUN apt-get update && \
    apt-get install --no-install-recommends -y procps && \
    rm -rf /var/lib/apt/lists/*

COPY actions/buildpack-aliyun/1.0/bp /opt/action/bp
COPY --from=docker-cli /usr/local/bin/docker /usr/bin/docker
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
COPY --from=builder /assets /opt/action

# $ tree /opt/action/
# |-- assets
# |   `-- java-agent
# |       |-- 3.10
# |       |   `-- spot-agent.tar.gz
# |       `-- 3.11
# |           `-- spot-agent.tar.gz
# |-- bp
# |   |-- dockerfile
# |   |-- java
# |   |   |-- build
# |   |   |   |-- maven
# |   |   |   |   |-- Dockerfile
# |   |   |   |   `-- maven-settings.xml
# |   |   |   `-- maven-edas
# |   |   |       |-- Dockerfile
# |   |   |       `-- rewrite_pom.py
# |   |   `-- pack
# |   |       |-- edas
# |   |       |   |-- Dockerfile
# |   |       |   `-- start.sh
# |   |       `-- springboot
# |   |           |-- Dockerfile
# |   |           `-- entrypoint.sh
# |   `-- node
# |       |-- build
# |       |   `-- npm
# |       |       `-- Dockerfile
# |       |
# |       `-- pack
# |           |-- herd
# |           |   |-- Dockerfile
# |           |   `-- entrypoint.sh
# |           `-- spa
# |               `-- Dockerfile
# `-- run
