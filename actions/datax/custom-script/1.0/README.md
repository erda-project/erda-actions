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
