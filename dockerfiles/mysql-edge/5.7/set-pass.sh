#!/bin/bash
password=${MYSQL_ROOT_PASSWORD:-R3M#ESuli0%aGj5v}
initpassword=$(cat /root/.mysql_secret | tail -n1)
while :;do
    mysqladmin -u root -p"$initpassword" password "$password" &> /dev/null && break
    sleep 10
done
rm -rf $0