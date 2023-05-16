FROM registry.erda.cloud/erda-actions/custom-script:2.0 AS builder

ENV CGO_ENABLED 0

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN go build -o /opt/action/run github.com/erda-project/erda-actions/actions/dockerfile/1.0/internal/cmd

FROM registry.erda.cloud/erda-actions/custom-script:2.0

COPY --from=builder /opt/action/run /opt/action/run
RUN chmod +x /opt/action/run
