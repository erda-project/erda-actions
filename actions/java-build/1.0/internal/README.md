## maven action使用

在pipeline stages里添加maven action:

```
maven:
    name: java-maven
    params:
        workdir: ${git-checkout}
        options: "-am -pl <module> -Dmaven.test.skip=true" // maven options
        target: "modulename/target/xxx.jar" // 相对路径，maven package的jar位置
        service: "modulename" // 与 dice.yml 里 service名称一致，用于镜像匹配
        profile: "default" // spring.profile.active
```

## mvn-action打包流程

1. make clean mvn-action
2. docker build -t xxx .
3. docker push xxx

## 内部实现流程

1. 加载 action & pipeline 环境变量
2. 利用环境变量渲染Dockerfile & maven settings 模板
3. 准备打包前所需文件(Dockerfile & settings) 至 docker build 上下文目录
4. 收集项目pom.xml, 构建应用缓存镜像
5. 执行 maven package 打包，docker build 构建应用镜像
6. docker push 应用镜像至集群 registry