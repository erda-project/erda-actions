FROM registry.erda.cloud/erda-x/golang:1.17 as builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN go build -o /assets/run           actions/dice/1.0/internal/run/cmd/*.go
RUN go build -o /assets/when_sigterm  actions/dice/1.0/internal/post/cmd/*.go

FROM registry.erda.cloud/erda-x/oraclelinux:7
COPY --from=builder /assets /opt/action
