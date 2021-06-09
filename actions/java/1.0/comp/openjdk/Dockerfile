# 可选: v1.8.0.242
ARG CONTAINER_VERSION=v1.8.0.242

FROM {{CENTRAL_REGISTRY}}/erda/terminus-openjdk:$CONTAINER_VERSION

ARG TARGET
ARG MONITOR_AGENT=true
ARG SCRIPT_ARGS

ENV SCRIPT_ARGS ${SCRIPT_ARGS}

RUN ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone

COPY comp/openjdk/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY pre_start.sh /pre_start.sh
RUN chmod +x /pre_start.sh

COPY comp/fonts /usr/share/fonts/custom
COPY comp/arthas-boot.jar /

ARG DICE_VERSION
COPY comp/spot-agent/${DICE_VERSION}/spot-agent.tar.gz /tmp/spot-agent.tar.gz
RUN \
	if [ "${MONITOR_AGENT}" = true ]; then \
        mkdir -p /opt/spot; tar -xzf /tmp/spot-agent.tar.gz -C /opt/spot; \
	fi && rm -rf /tmp/spot-agent.tar.gz

# copy target
RUN mkdir -p /app
COPY ${TARGET} /app/

COPY assets /assets

ENTRYPOINT ["/entrypoint.sh"]
