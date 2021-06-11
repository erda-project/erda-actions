FROM registry.erda.cloud/erda/terminus-golang:1.14 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /assets/run           actions/dice-deploy-rollback/1.0/internal/run/cmd/*.go
RUN GOOS=linux GOARCH=amd64 go build -o /assets/when_sigterm  actions/dice-deploy-rollback/1.0/internal/post/cmd/*.go

FROM registry.erda.cloud/erda/terminus-centos:base
COPY --from=builder /assets /opt/action
