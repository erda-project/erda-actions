FROM registry.erda.cloud/retag/pyroscope-java:v0.11.5 AS pyroscope-java
FROM registry.erda.cloud/erda-x/openjdk:8_11_17_21-ubuntu24.04-v1

ARG CONTAINER_VERSION=8
ENV CONTAINER_VERSION ${CONTAINER_VERSION}

ARG TARGET
ARG MONITOR_AGENT=true
ARG SCRIPT_ARGS

ENV SCRIPT_ARGS ${SCRIPT_ARGS}

COPY comp/openjdk/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY pre_start.sh /pre_start.sh
RUN chmod +x /pre_start.sh

COPY comp/fonts /usr/share/fonts/custom
#COPY comp/arthas-boot.jar /
COPY comp/jacocoagent.jar /opt/jacoco/jacocoagent.jar

ARG ERDA_VERSION
COPY comp/spot-agent/${ERDA_VERSION}/spot-agent-jdk17.tar.gz /tmp/spot-agent-jdk17.tar.gz
RUN \
	if [ "${MONITOR_AGENT}" = true ]; then \
        mkdir -p /opt/spot; tar --no-same-owner --no-same-permissions -xzf /tmp/spot-agent-jdk17.tar.gz -C /opt/spot; \
	fi && mv /opt/spot/spot-agent /opt/spot/spot-agent-jdk17 && rm -rf /tmp/spot-agent-jdk17.tar.gz

COPY comp/spot-agent/${ERDA_VERSION}/spot-agent.tar.gz /tmp/spot-agent.tar.gz
RUN \
	if [ "${MONITOR_AGENT}" = true ]; then \
        mkdir -p /opt/spot; tar --no-same-owner --no-same-permissions -xzf /tmp/spot-agent.tar.gz -C /opt/spot; \
	fi && rm -rf /tmp/spot-agent.tar.gz

COPY --from=pyroscope-java /app /opt/pyroscope

RUN mkdir -p /app
COPY ${TARGET} /app/

COPY assets /assets

ENTRYPOINT ["/entrypoint.sh"]
