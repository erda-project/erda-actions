FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-actions/custom-script:2.0 as builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN mkdir -p /opt/action/comp && cp -r actions/golang/1.0/comp/* /opt/action/comp

RUN go build -o /assets/run /go/src/github.com/erda-project/erda-actions/actions/golang/1.0/internal/cmd/main.go

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-actions/custom-script:2.0

COPY --from=builder /opt/action/comp /opt/action/comp
COPY --from=builder /assets /opt/action
