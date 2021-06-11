FROM registry.erda.cloud/erda/terminus-centos:base

WORKDIR /opt/action

RUN curl http://gosspublic.alicdn.com/ossutil/1.6.3/ossutil64  >  /bin/ossutil && chmod +x /bin/ossutil 

COPY actions/oss-upload/1.0/internal/run.sh /opt/action/run

RUN chmod +x /opt/action/run