# 可选: v8.5.43-jdk8, v7.0.96-jdk8
ARG CONTAINER_VERSION=v8.5.43-jdk8

FROM {{CENTRAL_REGISTRY}}/erda/terminus-tomcat:$CONTAINER_VERSION

ARG TARGET
ARG MONITOR_AGENT=true
ARG SCRIPT_ARGS
ARG WEB_PATH

ENV SCRIPT_ARGS ${SCRIPT_ARGS}

RUN ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone

COPY comp/tomcat/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY pre_start.sh /pre_start.sh
RUN chmod +x /pre_start.sh

ARG DICE_VERSION
COPY comp/spot-agent/${DICE_VERSION}/spot-agent.tar.gz /tmp/spot-agent.tar.gz
RUN \
	if [ "${MONITOR_AGENT}" = true ]; then \
        mkdir -p /opt/spot; tar -xzf /tmp/spot-agent.tar.gz -C /opt/spot; \
	fi && rm -rf /tmp/spot-agent.tar.gz

WORKDIR /

RUN cd ${CATALINA_HOME}/webapps/ && rm -fr *

# biz war
COPY ${TARGET}/app.war ${CATALINA_HOME}/webapps/${WEB_PATH}.war

CMD ["/entrypoint.sh"]
