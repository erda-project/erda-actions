FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-actions/custom-script:2.0 AS builder

COPY . /go/src/github.com/erda-project/erda-actions
WORKDIR /go/src/github.com/erda-project/erda-actions

RUN go build -o /opt/action/run github.com/erda-project/erda-actions/actions/java/1.0/internal/cmd

RUN mkdir -p /opt/action/comp && \
    cp -r actions/java/1.0/comp/* /opt/action/comp

# newest spot agent
RUN bash /opt/action/comp/download_spot_agent.sh
RUN bash /opt/action/comp/download_fonts.sh

FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda-actions/custom-script:2.0

COPY --from=builder /opt/action/run /opt/action/run
COPY --from=builder /opt/action/comp /opt/action/comp
