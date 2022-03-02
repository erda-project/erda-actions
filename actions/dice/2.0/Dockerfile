FROM registry.erda.cloud/retag/golang:1.16-alpine3.14 AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

# go build
RUN GOOS=linux GOARCH=amd64 go build -o /assets/run           actions/dice/2.0/internal/cmd/dice/*.go
RUN GOOS=linux GOARCH=amd64 go build -o /assets/when_sigterm  actions/dice/2.0/internal/cmd/cancel/*.go

FROM registry.erda.cloud/retag/alpine:3.14
COPY --from=builder /assets /opt/action
