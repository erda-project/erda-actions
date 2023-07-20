FROM registry.erda.cloud/retag/gcrane:v0.15.2 as tools
FROM registry.erda.cloud/erda/terminus-golang:1.19.7 as builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ENV GOPROXY=https://goproxy.cn,direct

RUN go build -o /opt/action/run github.com/erda-project/erda-actions/actions/docker-push/1.0/internal

FROM registry.erda.cloud/erda-actions/custom-script:2.0

COPY --from=builder /opt/action/run /opt/action/run
COPY --from=tools /ko-app/gcrane /usr/local/bin/gcrane

