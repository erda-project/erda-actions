# 1.0已弃用，请使用2.0，建议不要指定版本

### Custom-Script Action


运行自定义命令

#### 使用

```yaml
- custom-script:
    # 选填。声明自定义运行时镜像。
    # 平台默认提供的镜像包括 java, nodejs, golang 等编译环境。
    image: centos:7
    # 必填，执行的脚本命令列表，顺序执行。
    commands:
    - echo "hello world"
    - cat /etc/hosts
```

## Installed Software
### Language and Runtime
- Bash 5.0.17(1)-release
- Golang 1.17.12
- Node 12.22.12
- Python 2.7.18
- Python3 3.8.10

### Package Management
- Npm 6.14.16
- Pip 20.1.1
- Yarn 1.22.4

### Project Management
- Maven 3.6.3

### Tools
- Docker 20.10.3
- Git 2.26.3
- buildctl 0.9.2

### Java
| Version              | Vendor          | Environment Variable |
| -------------------- | --------------- | -------------------- |
| 8.275.01-r0          | Eclipse Temurin | JAVA_HOME_8_X64      |


```
    PostgreSQL service is disabled by default. Use the following command as a part of your job to start the service: 'sudo systemctl start postgresql.service'
```

### Cached Tools
#### Go
- 1.19.1

#### Node.js
- 12.22.12

#### Python
- 3.8.10

#### Environment variables
| Name            | Value                              | Architecture |
| --------------- | ---------------------------------- | ------------ |
| GOROOT | /usr/local/go | x64          |


### Installed packages
| Name           | Version                  |
|----------------|--------------------------|
| curl           | 7.79.1                   |
| iproute        | 1.34.1                   |
| libgit2-dev    | 1:1.2.11.dfsg-2ubuntu1.3 |
| tar            | 1.34.1                   |
| ssh | 1:8.2                    |
| scp |     1:8.2                     |
