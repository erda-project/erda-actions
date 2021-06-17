FROM registry.erda.cloud/erda/terminus-centos:7
COPY ["entrypoint.sh","install.sh","set-pass.sh","set-replica.sh","/usr/bin/"]
COPY ["master.cnf","slave.cnf", "/etc/mysql/"]
RUN ["/usr/bin/install.sh"]
VOLUME ["/var/lib/mysql"]
EXPOSE 3306
ENTRYPOINT ["/usr/bin/entrypoint.sh"]
#ENTRYPOINT ["sleep", "10000000"]