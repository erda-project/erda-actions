# syntax = docker/dockerfile:1.2

# Copyright (c) 2021 Terminus, Inc.
#
# This program is free software: you can use, redistribute, and/or modify
# it under the terms of the GNU Affero General Public License, version 3
# or later ("AGPL"), as published by the Free Software Foundation.
#
# This program is distributed in the hope that it will be useful, but WITHOUT
# ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
# FITNESS FOR A PARTICULAR PURPOSE.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>.

FROM registry.erda.cloud/retag/golang:1.15 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build action cmd
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    GOPROXY=https://goproxy.io,direct GOOS=linux GOARCH=amd64 go build -o /opt/action/run actions/erda-mysql-migration/1.0-57/internal/action/cmd/main.go

# go build local cmd
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    GOPROXY=https://goproxy.io,direct GOOS=linux GOARCH=amd64 go build -o /opt/action/erda-migrate actions/erda-mysql-migration/1.0-57/internal/local/cmd/main.go

RUN chmod 777 /opt/action/run
RUN chmod 777 /opt/action/erda-migrate

FROM registry.erda.cloud/retag/erda-mysql-migration-sandbox56:20210922

MAINTAINER chenzhongrun "zhongrun.czr@alibaba-inc.com"

USER root

VOLUME ["/log", "/migrations"]
ENV MIGRATION_DIR=/migrations

COPY --from=builder /opt/action/run /opt/action/run
COPY --from=builder /opt/action/erda-migrate /usr/bin/erda-migrate

RUN chmod 777 /opt/action/*
RUN chmod 777 /usr/bin/erda-migrate

ENTRYPOINT ["container-entrypoint"]
CMD ["erda-migrate"]
