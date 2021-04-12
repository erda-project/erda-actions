# registry.cn-hangzhou.aliyuncs.com/terminus/terminus-openjdk:v1.8.0.242-asia
FROM registry.cn-hangzhou.aliyuncs.com/terminus/terminus-openjdk:v1.8.0.242

RUN yum install -y unzip zip

RUN mkdir -p /asia/dice_files

RUN cd /asia/dice_files && wget http://arms-apm-shanghai.oss-cn-shanghai.aliyuncs.com/ArmsAgent.zip

RUN unzip /asia/dice_files/ArmsAgent.zip -d /asia/dice_files/
