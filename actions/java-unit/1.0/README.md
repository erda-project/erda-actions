### Java Action

用于 java 单元测试，支持 JUnit 和 TestNg 框架，制作成为 docker image 用于运行服务。

#### 使用

Examples:
1. pipeline.yml 增加 lint 描述
```yml
- java-unit:
    params:
      path: xxxx
```
