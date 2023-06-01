FROM registry.erda.cloud/retag/pyroscope-java:v0.11.3 as pyroscope-java
FROM registry.erda.cloud/erda-x/tomcat:8.5

ARG CONTAINER_VERSION=v8
ENV CONTAINER_VERSION ${CONTAINER_VERSION}

ARG TARGET
ARG MONITOR_AGENT=true
ARG SCRIPT_ARGS
ARG WEB_PATH

ENV SCRIPT_ARGS ${SCRIPT_ARGS}

COPY comp/tomcat/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY pre_start.sh /pre_start.sh
RUN chmod +x /pre_start.sh

ARG ERDA_VERSION
COPY comp/spot-agent/${ERDA_VERSION}/spot-agent.tar.gz /tmp/spot-agent.tar.gz
RUN \
	if [ "${MONITOR_AGENT}" = true ]; then \
        mkdir -p /opt/spot; tar -xzf /tmp/spot-agent.tar.gz -C /opt/spot; \
	fi && rm -rf /tmp/spot-agent.tar.gz

COPY --from=pyroscope-java /app /opt/pyroscope
COPY comp/jacocoagent.jar /opt/jacoco/jacocoagent.jar

WORKDIR /

RUN cd ${CATALINA_HOME}/webapps/ && rm -fr *

COPY ${TARGET}/app.war ${CATALINA_HOME}/webapps/${WEB_PATH}.war

CMD ["/entrypoint.sh"]
