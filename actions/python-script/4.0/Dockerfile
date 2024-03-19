FROM --platform=$BUILDPLATFORM ubuntu:latest

# The original 4.0 version is missing. This is a simulation supplement and can be used normally
MAINTAINER jingzhi jingzhi736728@163.com>

# Installation dependencies
RUN apt-get update -y && apt-get install -y python3

# Installing Python related libraries
RUN apt-get -y install python3-pip
RUN python3 -m pip install aliyun-python-sdk-core
RUN python3 -m pip install oss2
RUN python3 -m pip install kafka-python
RUN python3 -m pip install redis
RUN python3 -m pip install elasticsearch
RUN python3 -m pip install thrift
RUN python3 -m pip install pyhive

# Set up soft links
RUN ln -s /usr/bin/python3 /usr/bin/python

# working directory
WORKDIR /

CMD ["/bin/bash"]