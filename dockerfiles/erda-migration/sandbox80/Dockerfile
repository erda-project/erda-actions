FROM registry.erda.cloud/erda-addons/mylet:v8.0

USER root:root
ENTRYPOINT ["docker-entrypoint.sh"]
EXPOSE 3306 33060
CMD ["mysqld"]

ENV MYSQL_ROOT_PASSWORD="12345678"
ENV SANDBOX_INNER_PASSWORD="12345678"
ENV MYSQL_ALLOW_EMPTY_PASSWORD yes

COPY dockerfiles/erda-migration/my.cnf /etc/my.cnf

RUN yum erase -y mysql-community-server-minimal && yum install -y mysql-server mysql-devel gcc python3 python3-devel

# https://bugs.mysql.com/bug.php?id=91395
RUN setcap -r /usr/libexec/mysqld

RUN python3 -m pip --no-cache-dir install pip -U
RUN python3 -m pip --no-cache-dir config set global.index-url https://mirrors.aliyun.com/pypi/simple
RUN python3 -m pip --no-cache-dir install mysql-connector-python==8.0.30 django==3.2.4 pytz==2021.1 sqlparse==0.4.1

RUN mkdir /logs
RUN touch /logs/healthy
