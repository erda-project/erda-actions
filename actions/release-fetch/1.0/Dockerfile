FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-x/golang:1.22 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

# go build
RUN GOOS=linux GOARCH=$ARCH GOPROXY=https://goproxy.cn go build -o /assets/run actions/release-fetch/1.0/main.go

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-x/debian-bookworm:12
COPY --from=builder /assets /opt/action
