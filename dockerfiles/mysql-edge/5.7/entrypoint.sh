#!/bin/bash

# mysql 日志目录、数据目录(数据目录根据节点数据盘具体调整)
mkdir -p /var/log/mysql
touch /var/log/mysql/error.log
mkdir -p /data
chmod 0777 /tmp

echo "execute mysql install db"
[ -d "/data/mysql" ] || mysql_install_db --user=mysql --basedir=/usr --datadir=/data/mysql

# 设置 目录用户组(数据目录根据节点数据盘具体调整)
chown -R mysql:mysql /var/log/mysql
chown -R mysql:mysql /data/mysql

(nohup set-pass.sh &)
(nohup set-replica.sh &)
# 如果要配置AB复制集群，根据Headless Service，主服务器的 podname 一定是以 “-0” 结尾
# 主从服务器配置文件不一样，启动前要配置不同的配置文件。
# 务必要用ConfigMap传不同的配置文件（master.cnf,slave.cnf）到/etc/mysql目录下。
if [ -n "$POD_NAME" ];then
    pod_seq=$(awk -F"[-]" '{print $2}' <<< $POD_NAME)
    server_id=$[pod_seq+1]
    if [[ $POD_NAME =~ "mysql-master" ]];then
        sed -r "s/(server.id).*/\1 = 1/" /etc/mysql/master.cnf > /etc/my.cnf
    else
        sed -r "s/(server.id).*/\1 = 2/" /etc/mysql/slave.cnf > /etc/my.cnf
    fi
fi
# 临时
#cp  /etc/mysql/master.cnf  /etc/my.cnf
#临时 end
echo "starting mysqld"
#/usr/sbin/mysqld --pid-file=/var/run/mysqld/mysqld.pid --user=mysql
(tail -f /var/log/mysql/error.log  > /dev/stdout &)
/usr/sbin/mysqld