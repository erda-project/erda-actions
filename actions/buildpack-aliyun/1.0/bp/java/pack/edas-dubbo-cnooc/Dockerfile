# Build app image
FROM --platform=$TARGETPLATFORM registry.erda.cloud/erda/edas-container:3.4.4-cnooc

RUN echo "Asia/Shanghai" | tee /etc/timezone

# set default TZ, modify through `--build-arg TZ=XXX`
ARG TZ="Asia/Shanghai"

WORKDIR /
RUN wget https://arms-apm-cn-hangzhou.oss-cn-hangzhou.aliyuncs.com/2.7.1.4/ArmsAgent.zip -O /opt/ArmsAgent.zip
RUN rm -rf /cnooc/dice_files/ArmsAgent && unzip /opt/ArmsAgent.zip -d /cnooc/dice_files/

COPY /bp/pack/start.sh /home/admin/bin/start.sh
RUN chmod +x /home/admin/bin/start.sh

ARG USE_AGENT=true
ARG DICE_VERSION
COPY /assets/java-agent/${DICE_VERSION}/spot-agent.tar.gz /tmp/spot-agent.tgz

RUN \
    if [ $USE_AGENT = true ]; then \
        mkdir -p /opt/spot; tar -xzf /tmp/spot-agent.tgz -C /opt/spot; rm -rf /tmp/spot-agent.tgz; \
    fi

ENV SPRING_PROFILES_ACTIVE=default
RUN yum -y install nc

ADD /app /home/admin/app
