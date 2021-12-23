FROM registry.erda.cloud/retag/buildkit:v0.9.2 as buildkit

FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /assets/run github.com/erda-project/erda-actions/actions/buildpack/1.0/internal

# async assets
WORKDIR actions/buildpack/1.0
RUN bash sync_assets.sh

FROM registry.erda.cloud/erda/centos:7
RUN yum install -y yum-utils && \
    yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo && \
    yum -y install docker-ce-cli

COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl

COPY actions/buildpack/1.0/bp /opt/action/bp
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
