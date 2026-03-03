FROM registry.erda.cloud/erda-x/golang:1.24 AS builder

ENV CGO_ENABLED=0

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN mkdir -p /opt/action/comp && cp -r actions/golang/1.0/comp/* /opt/action/comp

RUN go build -o /assets/run /go/src/github.com/erda-project/erda-actions/actions/golang/1.0/internal/cmd/main.go

FROM registry.erda.cloud/retag/buildkit:v0.11.6 AS buildkit
FROM registry.erda.cloud/retag/docker:20.10.24-cli AS docker-cli
FROM registry.erda.cloud/erda-x/golang:1.24

COPY --from=builder /opt/action/comp /opt/action/comp
COPY --from=builder /assets /opt/action
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
COPY --from=docker-cli /usr/local/bin/docker /usr/bin/docker

RUN chmod +x /opt/action/run
