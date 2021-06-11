### SQL Action

### 一、简介

支持与sftp server交互的工具，目前只支持从sftp server获取文件。
### 二、详细介绍
用户只需要配置一些必要的参数，即可从sftp server获取文件。
配置参数含义如下：
SFTP_HOST:  sftp server的主机地址
SFTP_PORT:  sftp server的端口号
SFTP_USER:  登录使用的用户名
SFTP_PASSWORD:  登录使用的密码
SFTP_REMOTE_PATH:   需要下载的远程文件
SFTP_LOCAL_PATH:    本地文件保存路径

### 三、使用场景
用户需要通过sftp协议去下载文件。

### 四、使用方式
```aidl
version: '1.0'

resource_types:
  - name: sftp
    type: docker-image
    source:
      repository: registry.erda.cloud/erda/pipeline-sftp-resource
      tag: v1.0.0

resources:
  - name: sftp-get
    type: sftp
    source:
      SFTP_HOST: 1.1.1.1
      SFTP_PORT: 22
      SFTP_USER: foo
      SFTP_PASSWORD: foo

stages:
- name: sftp-get
  tasks:
    - put: sftp-get
      params:
        SFTP_REMOTE_PATH: /root/y.json
        SFTP_LOCAL_PATH: y.json
```
