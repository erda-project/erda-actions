#!/bin/bash
cd /etc/yum.repos.d
cat > mysql-57.repo << eof
[mysql57-community]
name=MySQL 5.7 Community Server
baseurl=http://repo.mysql.com/yum/mysql-5.7-community/el/7/x86_64/
enabled=1
gpgcheck=0
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-mysql
eof
yum install mysql-community-server -y
yum clean all
rm -rf $0