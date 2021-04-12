FROM busybox
LABEL VERSION=3.20
RUN mkdir -p /resource/java-agent
RUN mkdir -p /dice
RUN wget https://terminus-dice.oss-cn-hangzhou.aliyuncs.com/spot/java-agent/action/release/3.20/spot-agent.tar.gz -O /resource/java-agent/spot-agent.tar.gz

ENTRYPOINT tar -xf /resource/java-agent/spot-agent.tar.gz -C /dice