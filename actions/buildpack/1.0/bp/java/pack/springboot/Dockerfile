# source: dockerfiles/pyroscope/java/Dockerfile
# /app/pyroscope.jar natively supports amd64/arm64 arch.
FROM registry.erda.cloud/retag/pyroscope-java:v0.11.3 as pyroscope-java
# Build app image
FROM registry.erda.cloud/erda-x/openjdk:8_11

# set default TZ, modify through `--build-arg TZ=XXX`
ARG TZ="Asia/Shanghai"

RUN sed -i 's/jdk.tls.disabledAlgorithms=SSLv3, TLSv1, TLSv1.1/jdk.tls.disabledAlgorithms=SSLv3/g' /usr/lib/jvm/java-1.8.0-openjdk/jre/lib/security/java.security

COPY /bp/pack/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

WORKDIR /

ARG USE_AGENT=true
ARG USE_PREVIEW_AGENT=false
ARG DICE_VERSION
COPY /assets/java-agent/${DICE_VERSION}/spot-agent.tar.gz /tmp/spot-agent.tgz
#COPY /bp/pack/arthas-boot.jar /
COPY /bp/pack/jacocoagent.jar /opt/jacoco/jacocoagent.jar

RUN \
    if [ $USE_AGENT = true ]; then \
        mkdir -p /opt/spot; tar -xzf /tmp/spot-agent.tgz -C /opt/spot; rm -rf /tmp/spot-agent.tgz; \
    fi

COPY --from=pyroscope-java /app /opt/pyroscope

ENV SPRING_PROFILES_ACTIVE=default

ADD /app /app

CMD ["/entrypoint.sh"]

