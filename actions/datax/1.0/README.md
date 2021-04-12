### DataX Action

### 一、简介

DataX为用户提供多种不同数据源之间的全量同步功能。

### 二、详细介绍

DataX是由阿里巴巴集团开源的数据同步工具。DataX支持常见的数据存储系统间的数据全量同步。DataX action增加了DataX支持的数据源种类，同时简化了参数配置，提高了易用性。

### 三、使用场景
DataX的使用场景是异构数据源之间离线同步数据，目前支持的数据源有：MySQL、Oracle、DB2、HDFS、Hive、ODPS、HBase、Cassandra、FTP等。

### 四、使用方式
```aidl
version: '1.0'
resources:
  - name: repo
    type: git
    source:
      uri: ((gittar.repo))
      branch: ((gittar.branch))
      username: ((gittar.username))
      password: ((gittar.password))
  - name: datax
    type: datax
    source:
        jsonFilePath: "repo/ea/etl_terminus/ods/source/import_commit_detail_gittar_init.json"
        databaseType: hive
        url: "jdbc:hive2://1.1.1.1:9000/;auth=noSasl"
        username: "foo"
        password: "foo"
        database: "git_accompany"
  - name: s_commit_detail-init
    type: sql
    source:
        queryType: sparksql
        queryEndPoint: jdbc:hive2://1.1.1.1:9000/;auth=noSasl
        username: foo
        password: foo
stages:
  - name: repo
    tasks:
      - get: repo
        params:
          depth: 3
  - name: datax
    tasks:
        - put: datax
          params:
            outputTables:
              - s_commit_detail_gittar_init
  - name: s_commit_detail-init
    tasks:
        - put: s_commit_detail-init
          params:
            path: repo/ea/etl_terminus/ods/init/s_commit_detail_gittar_init.q
            outputTables:
              - git.s_commit_detail_gittar
            process: init
```
