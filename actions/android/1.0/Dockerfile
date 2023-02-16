FROM registry.erda.cloud/erda/terminus-golang:1.17.12 as builder

MAINTAINER shenli shenli@terminus.io

ENV GOLANG_VERSION 1.17.12

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
ENV GOPROXY=https://goproxy.io,https://goproxy.cn,direct
RUN GOOS=linux GOARCH=amd64 go build -o /assets/run /go/src/github.com/erda-project/erda-actions/actions/android/1.0/internal/cmd/main.go

FROM registry.erda.cloud/erda/android-gradle-node:v30-20220721203530
ENV SASS_BINARY_SITE="https://npmmirror.com/mirrors/node-sass"
ENV NODEJS_ORG_MIRROR="https://npmmirror.com/dist"
# built-in installation of some google sdk
RUN sdkmanager "build-tools;28.0.3" "build-tools;31.0.0"
RUN sdkmanager "platforms;android-28" "platforms;android-31" "platforms;android-30"
RUN npm config set @terminus:registry http://registry.npm.terminus.io/ && \
    npm config set registry http://registry.npmmirror.com/
RUN apt-get update
RUN apt-get install -y openjdk-11-jdk
COPY actions/android/1.0/comp/init.gradle /root/.gradle/init.gradle
COPY --from=builder /assets /opt/action
