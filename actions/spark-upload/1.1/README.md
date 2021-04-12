### Spark-Upload Action

### 一、简介

Spark上传 Action 允许用户上传Spark应用的jar包。

### 二、详细介绍

Spark是高性能的基于内存计算的分布式大数据处理框架。Spark上传 Action允许用户上传Spark应用jar包。使用该Action，用户需要先将程序jar包上传至对象存储系统（Minio/OSS），然后生成一个http地址给下游Action使用。
### 三、使用场景
Spark-Upload使用场景是当需要以jar包的形式向一个Spark集群提交任务时，可以使用该action将jar包上传到对象存储以得到一个可以通过http方式访问jar包的方式。

### 四、使用方式
```aidl
type: spark-upload

desc: Spark上传 Action 允许用户上传Spark应用的jar包

support:
  get: false
  put: true

source:

  - name: jarPath
    required: true
    desc: 云存储jar包路径

params:

labels:
  maintainer: cb167668@alibaba-inc.com
```
