### Custom-Script Action

自定义命令2.0版本，在1.0的基础上对工具(golang,sed,mvn,node,curl,docker等)进行了升级

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
