FROM registry.erda.cloud/retag/buildkit:v0.9.2 as buildkit
FROM registry.erda.cloud/erda/terminus-golang:1.14 as builder

MAINTAINER shenli shenli@terminus.io

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions
RUN mkdir -p /opt/action/comp && \
    cp -r actions/php/1.0/comp/* /opt/action/comp

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /assets/run /go/src/github.com/erda-project/erda-actions/actions/php/1.0/internal/cmd/main.go

FROM registry.erda.cloud/erda/terminus-centos:base
COPY --from=composer /usr/bin/composer /usr/bin/composer
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
RUN rpm -Uvh https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm \
 && rpm -Uvh https://mirror.webtatic.com/yum/el7/webtatic-release.rpm
RUN yum install -y  php72w-cli php72w-mysql php72w-pdo php72w-xml php72w-mbstring php72w-gd
RUN yum install -y  docker make
RUN yum install -y unzip git
RUN composer config -g repo.packagist composer https://mirrors.aliyun.com/composer/
RUN composer global require slince/composer-registry-manager
RUN composer repo:use tencent
RUN composer config -g process-timeout 600
COPY --from=builder /assets /opt/action
COPY --from=builder /opt/action/comp /opt/action/comp
