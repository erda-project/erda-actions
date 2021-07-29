### Java-build Action

Warning!!! 当前为 Beta 版，不推荐使用，不向前兼容。

用于编译打包 java 工程, 推送 jar 包到私服仓库, 内置 maven 和 gradle

#### 使用

Examples:

1 maven

```yaml
- java-build:
    alias: java-build
    params:
      build_cmd:  # 用于构建的命令，可以多行
        - "mvn clean install -Dmaven.test.skip=true"
      jdk_version: "8" # 目前只支持 8 和 11
      workdir: ${git-checkout} # 在哪个目录下执行命令，一般是在代码目录下
```

2 gradle

```yaml
- java-build:
    alias: java-build
    params:
      build_cmd: 
        - "./gradlew build"
      jdk_version: "8"
      workdir: ${git-checkout}
```


#### 配合 release 的 services 使用


```yaml
java-demo:  # 声明 service 的名称，名称对应 dice,yaml 中的 service
    # 在哪个环境下运行，一般平台图形界面会给你设置好 image, 你只需要选择对应的 jdk 版本
    image: "openjdk:8-jre-alpine"
    copys:
      # 将当前 java-build 生成的 jar 包拷贝到当前 /target 下，这样就可以配合 java -jar 直接运行
      - ${java-build:OUTPUT:buildPath}/target/docker-java-app-example.jar:/target
      # 固定加上即可，用于配合 ${java-build:OUTPUT:JAVA_OPTS} 的监控
      - ${java-build:OUTPUT:buildPath}/spot-agent:/
     # 项目运行的命令, ${java-build:OUTPUT:JAVA_OPTS} 是对应的监控命令，固定加上即可
     cmd: java ${java-build:OUTPUT:JAVA_OPTS} -jar /target/docker-java-app-example.jar 
```

tomcat

可以用替换 ROOT.war 的方式在 tomcat 的根目录运行
export JAVA_OPTS="${java-build:OUTPUT:JAVA_OPTS}" \
&& mv /usr/local/tomcat/webapps/docker-java-app-example.war /usr/local/tomcat/webapps/ROOT.war \
 && /usr/local/tomcat/bin/catalina.sh run


```yaml
java-demo: 
    # 图形界面基于选择，当然可以自选
    image: "tomcat:jdk8-openjdk-slim"
    copys:
      # 拷贝 war 包到 webapps 下
      - ${java-build:OUTPUT:buildPath}/target/docker-java-app-example.war:/usr/local/tomcat/webapps
      # 设置固定值
      - ${java-build:OUTPUT:buildPath}/spot-agent:/
    # 启动 catalina.sh 并设置固定值
    cmd: export JAVA_OPTS="${java-build:OUTPUT:JAVA_OPTS}" && /usr/local/tomcat/bin/catalina.sh run
```
