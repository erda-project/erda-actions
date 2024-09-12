FROM --platform=$TARGETPLATFORM registry.erda.cloud/retag/buildkit:v0.11.3 as buildkit
FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-x/golang:1.22 AS builder

WORKDIR /go/src/github.com/erda-project/erda-actions

ENV CGO_ENABLED=0

COPY . .

RUN go build -o /opt/action/run github.com/erda-project/erda-actions/actions/dockerfile/1.0/internal/cmd

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-x/debian-bookworm:12

COPY --from=builder /opt/action/run /opt/action/run
COPY --from=buildkit /usr/bin/buildctl /usr/bin/buildctl
RUN chmod +x /opt/action/run
