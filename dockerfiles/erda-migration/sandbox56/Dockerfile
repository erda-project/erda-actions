FROM registry.erda.cloud/retag/mysql-56-centos7

MAINTAINER chenzhongrun "zhongrun.czr@alibaba-inc.com"

USER root

ENV TZ=Asia/Shanghai
ENV MYSQL_ROOT_PASSWORD="12345678"
ENV MYSQL_ALLOW_EMPTY_PASSWORD yes

COPY dockerfiles/erda-migration/sandbox56/my.cnf /etc/my.cnf

# install python3
RUN yum install -y mariadb-devel gcc python3-devel.x86_64
RUN python3 -m pip --no-cache-dir install pip -U
RUN python3 -m pip --no-cache-dir config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
RUN python3 -m pip --no-cache-dir install mysql-connector-python==8.0.26 django==3.2.4 pytz==2021.1 sqlparse==0.4.1

# support chinese
RUN yum install kde-l10n-Chinese -y
RUN yum install glibc-common -y
RUN localedef -c -f UTF-8 -i zh_CN zh_CN.utf8
ENV LC_ALL zh_CN.UTF-8

# add healthy probe
RUN mkdir /logs
RUN touch /logs/healthy
