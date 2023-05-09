FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-actions/custom-script:2.0 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN go build -o /assets/run github.com/erda-project/erda-actions/actions/buildpack/1.0/internal

WORKDIR actions/buildpack/1.0
RUN bash sync_assets.sh

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-actions/custom-script:2.0

COPY --from=builder /assets /opt/action
COPY actions/buildpack/1.0/bp /opt/action/bp
