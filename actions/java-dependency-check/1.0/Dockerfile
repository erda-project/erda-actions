FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

LABEL maintainer="jiangzhengdong <zhengdong.jzd@alibaba-inc.com>"

# disable CGO for ALL THE THINGS (to help ensure no libc)
ENV CGO_ENABLED 0
ENV BUILD_FLAGS="-v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo"

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

COPY actions/java-dependency-check/1.0/settings.xml /opt/action/mvn/settings.xml

RUN set -x \
    	&& eval "GOOS=linux GOARCH=amd64 go build $BUILD_FLAGS -o /opt/action/run github.com/erda-project/erda-actions/actions/java-dependency-check/1.0/internal/cmd"

# 从 dockerhub 每天更新的 dcdb 镜像中获取最新的漏洞数据库
FROM adorsys/dependency-check-db:6.0-h2 AS dcdb

FROM registry.erda.cloud/erda/terminus-maven:3-jdk-8-alpine

COPY --from=builder /opt/action/run /opt/action/run
COPY --from=builder /opt/action/mvn/settings.xml /opt/action/mvn/settings.xml
COPY --from=dcdb /usr/share/nginx/html /opt/action/dependency-check
RUN mvn org.owasp:dependency-check-maven:6.3.1:update-only -DdataDirectory=/opt/action/dependency-check


CMD ["/opt/action/run"]
