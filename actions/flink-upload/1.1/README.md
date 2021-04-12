### Flink-Upload Action

### 一、简介

Flink上传 Action允许用户上传一个Flink作业的jar包

### 二、详细介绍

Flink 是一种高效的流式计算引擎，同时也支持批处理作业，正在朝着流批一体的方向演进。Flink上传 Action 允许用户将应用jar包上传到指定的Flink集群上。为了使用该Action,用户需要提供Flink Jobmanager地址和应用jar包路径。

### 三、使用场景

Flink-Upload的使用场景是：向Flink集群提交作业时，通过该action将jar包上传到JobManager管理的路径。

### 四、使用方式
```aidl
version: '1.0'

resource_types:
- name: flink-upload
  type: docker-image
  source:
    repository: bertonchen/pipeline-flink-resource
    tag: v2.8
- name: file-download
  type: docker-image
  source:
    repository: bertonchen/pipeline-cloudstorage-resource
    tag: v4.7

stages:
- name: get-file
  tasks:
  - put: file
    params:
      srcFile: "WordCount.jar"
      destFile: "WordCount.jar"
- name: flink-upload
  tasks:
  - put: flink-upload
- name: flink
  tasks:
  - put: flink

resources:
- name: file
  type: file-download
  source:
    cloudType: "MINIO"
    operateType: "download"
    endpoint: "http://1.1.1.1:9000"
    bucketName: "foo"
    accessKey: "foo"
    secretKey: "foo"
- name: flink-upload
  type: flink-upload
  source:
    jobManagerUrl: "http://1.1.1.1:8081"
    jarPath: "file/WordCount.jar"
- name: flink
  type: flink
  source:
    labels:
      JOB_KIND: bigdata
    depends: flink-upload
```
