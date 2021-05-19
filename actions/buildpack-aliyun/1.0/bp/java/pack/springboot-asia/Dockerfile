# Build app image
FROM {{BP_DOCKER_BASE_REGISTRY}}/erda/terminus-openjdk:v1.8.0.242-asia

# set default TZ, modify through `--build-arg TZ=XXX`
ARG TZ="Asia/Shanghai"

COPY /bp/pack/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

WORKDIR /

ARG USE_AGENT=true
ARG USE_PREVIEW_AGENT=false
ARG DICE_VERSION
COPY /assets/java-agent/${DICE_VERSION}/spot-agent.tar.gz /tmp/spot-agent.tgz
COPY /bp/pack/arthas-boot.jar /

RUN \
    if [ $USE_AGENT = true ]; then \
        mkdir -p /opt/spot; tar -xzf /tmp/spot-agent.tgz -C /opt/spot; rm -rf /tmp/spot-agent.tgz; \
    fi

ENV SPRING_PROFILES_ACTIVE=default

ADD /app /app

CMD ["/entrypoint.sh"]

