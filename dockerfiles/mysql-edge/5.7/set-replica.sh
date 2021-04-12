#!/bin/bash
# 如果是AB复制集群，则需要通过downwardAPI获取podname，并赋值给POD_NAME变量
# 如果没有进行上面的步骤，则说明是单节点mysql，本脚本直接退出
[ -z "$POD_NAME" ] && exit
password=${MYSQL_ROOT_PASSWORD:-R3M#ESuli0%aGj5v}
# 当数据库服务启动成功后进行判断，结合 Headless Service，主服务器的 podname 一定是以 “-0” 结尾
# 主服务器执行 GRANT 操作，从服务器执行 CHANGE MASTER 操作
while :;do
    mysql -u root -p$password -e "select 1;" &> /dev/null  &&
    if [[ $POD_NAME =~ "mysql-master" ]];then
        mysql -u root -p$password -e "grant replication slave,replication client on *.* to slave@'%' identified by 'Bao_12345678';" > /dev/null
        mysql -u root -p$password -e "create user 'mysql'@'%' identified by '$MYSQL_USER_PASSWORD';" > /dev/null
        mysql -u root -p$password -e "grant all privileges on *.* to 'mysql'@'%' WITH GRANT OPTION;" > /dev/null
        mysql -u root -p$password -e "flush privileges;" > /dev/null
    else
        mysql -u root -p$password -e "change master to \
                                      master_host='$MASTER_SVC_NAME', \
                                      master_user='slave', \
                                      master_password='Bao_12345678', \
                                      MASTER_AUTO_POSITION=1"
        mysql -u root -p$password -e "reset slave;" > /dev/null
        mysql -u root -p$password -e "start slave;" > /dev/null
    fi && break
    sleep 1
done