### File-sync Action

### 一、简介

文件同步Action允许用户从对象存储系统中上传或下载文件,当前支持Minio和阿里云OSS。
### 二、详细介绍

文件同步Action为用户提供文件上传或下载功能。用户只需要配置源文件所在地址 ，以及文件要上传或下载的位置，即可被使用。目前文件支持MINIO和OSS两种云存储的方式。

### 三、使用场景
File-sync的使用场景有两个：一是将文件上传到对象存储系统，而是从对象存储系统中下载文件到本地路径。

### 四、使用方式
```aidl
version: '1.0'
resources:
- name: file
  type: file-sync
  source:
    cloudType: "MINIO"
    operateType: "download"
    endpoint: "http://1.1.1.1:9000"
    bucketName: "horus"
    accessKey: "xxx"
    secretKey: "xxx"
stages:
- name: get-file
  tasks:
  - put: file
    params:
      srcFile: "resource/1/code-crawler-1.0-SNAPSHOT.jar"
      destFile: "code-crawler-1.0-SNAPSHOT.jar"
- name: exec
  tasks:
  - task: java-jar
    config:
      image_resource:
        type: docker-image
        source:
          repository: bertonchen/pipeline-cloudstorage-resource
      envs:
        TEST_ENV: test_env_task_value_3
      inputs:
      - name: file
      run:
        path: java
        args:
        - -jar
        - file/code-crawler-1.0-SNAPSHOT.jar
```
