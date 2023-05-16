FROM registry.erda.cloud/erda-x/golang:1.17 as builder

ENV CGO_ENABLED=0

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN go build -o /assets/run           actions/dice/2.0/internal/cmd/dice/*.go
RUN go build -o /assets/when_sigterm  actions/dice/2.0/internal/cmd/cancel/*.go

FROM registry.erda.cloud/erda-x/oraclelinux:7
COPY --from=builder /assets /opt/action
