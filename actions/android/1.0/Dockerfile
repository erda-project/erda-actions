FROM registry.erda.cloud/erda/terminus-golang:1.14 as builder

MAINTAINER shenli shenli@terminus.io

ENV GOLANG_VERSION 1.12.7

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /assets/run /go/src/github.com/erda-project/erda-actions/actions/android/1.0/internal/cmd/main.go

FROM registry.erda.cloud/erda/android-gradle-node:v29
ENV SASS_BINARY_SITE="https://npm.taobao.org/mirrors/node-sass"
ENV NODEJS_ORG_MIRROR="https://npm.taobao.org/dist"
RUN npm config set @terminus:registry http://registry.npm.terminus.io/ && \
    npm config set registry http://registry.npm.taobao.org/
COPY actions/android/1.0/comp/init.gradle /root/.gradle/init.gradle
COPY --from=builder /assets /opt/action