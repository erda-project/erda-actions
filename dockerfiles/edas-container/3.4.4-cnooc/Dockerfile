# registry.erda.cloud/erda/edas-container:3.4.4-cnooc
FROM registry.erda.cloud/erda/edas-container:3.4.4

RUN yum install -y unzip zip

RUN mkdir -p /cnooc/dice_files

RUN cd /cnooc/dice_files && wget http://terminus-dice.oss-cn-hangzhou.aliyuncs.com/spot/java-agent/action/release/cnooc/ArmsAgent.zip

RUN unzip /cnooc/dice_files/ArmsAgent.zip -d /cnooc/dice_files/
