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

FROM golang:1.15 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build action cmd
RUN GOPROXY=https://goproxy.io,direct GOOS=linux GOARCH=amd64 go build -o /opt/action/run actions/erda-mysql-migration/1.0/internal/action/cmd/main.go

# go build local cmd
RUN GOPROXY=https://goproxy.io,direct GOOS=linux GOARCH=amd64 go build -o /opt/action/erda-migrate actions/erda-mysql-migration/1.0/internal/local/cmd/main.go

RUN chmod 777 /opt/action/run
RUN chmod 777 /opt/action/erda-migrate

FROM centos/mysql-57-centos7

MAINTAINER chenzhongrun "zhongrun.czr@alibaba-inc.com"

USER root

ENV MYSQL_ROOT_PASSWORD="12345678"
ENV MYSQL_ALLOW_EMPTY_PASSWORD yes
ENV MIGRATION_DIR=/migrations

COPY --from=builder /opt/action/run /opt/action/run
COPY --from=builder /opt/action/erda-migrate /usr/bin/erda-migrate

# write mysql config
RUN bash -c "echo [mysqld]                    >  /etc/my.cnf"
RUN bash -c "echo user=mysql                  >> /etc/my.cnf"
RUN bash -c "echo symbolic-links = 0          >> /etc/my.cnf"
RUN bash -c "echo skip_name_resolve           >> /etc/my.cnf"
RUN bash -c "echo \!includedir /etc/my.cnf.d  >> /etc/my.cnf"

RUN mkdir -p /migrations
RUN chmod 777 /opt/action/*
RUN chmod 777 /usr/bin/erda-migrate
RUN chmod 777 /migrations

RUN yum install -y mariadb-devel gcc python3-devel.x86_64
RUN python3 -m pip --no-cache-dir install pip -U
RUN python3 -m pip --no-cache-dir config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
RUN python3 -m pip --no-cache-dir install mysqlclient==2.0.3 django==3.2.4 pytz==2021.1 sqlparse==0.4.1

VOLUME ["/migrations"]

ENTRYPOINT ["erda-migrate"]
