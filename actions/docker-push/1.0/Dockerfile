FROM --platform=$BUILDPLATFORM gcr.io/go-containerregistry/gcrane:v0.7.0 as gcrane
FROM --platform=$BUILDPLATFORM registry.erda.cloud/erda-x/terminus-golang:1.19.7

ARG TARGETARCH

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

ARG ARCH

RUN GOOS=linux GOARCH=$TARGETARCH go build -o /opt/action/run github.com/erda-project/erda-actions/actions/docker-push/1.0/internal

COPY --from=gcrane /ko-app/gcrane /usr/bin
