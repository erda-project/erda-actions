# Java Action

用于编译打包 java 工程，制作成为 docker image 用于运行服务。

## outputs

outputs 可以通过 `${alias:OUTPUT:output}` 的方式被后续 action 引用。

支持的 outputs 列表如下：

- image

  示例：

  ```yaml
  - stage:
    - java:
        ......
  - stage:
    - release:
        params:
          image: ${java:OUTPUT:image}
  ```

## 使用

Examples:

### maven / spring-boot

```yaml
- java:
    params:
      workdir: ${git-checkout} # 打包时的当前路径，此路径一般为根 pom.xml 的路径
      build_type: maven # 打包类型，这里选用 maven 打包
      options: -am -pl user # maven 打包参数，比如打包 user 模块使用命令 `mvn clean package -am -pl user`，这里省略命令 `mvn clean package` 只需要填写参数
      target: ./user/target/user.jar # 打包产物，一般为 jar，填写相较于 workdir 的相对路径。文件必须存在，否则将会出错。
      container_type: spring-boot # 运行 target（如 jar）所需的容器类型，比如这里我们打包的结果是 spring-boot 的 fat jar，故使用 spring-boot container
      #container_version: 8 # 默认和 jdk_version 一致。可选: 8 / 11 / 17 / 21
```

```yaml
- java:
    alias: java-demo
    params:
      jdk_version: 11 # 8 (默认) / 11 / 17 / 21
      build_type: maven
```

### maven / tomcat

```yaml
- java:
    params:
      jdk_version: 8
      build_type: maven
      workdir: ${git-checkout}
      options: -am -pl user
      target: ./user/target/user.war
      container_type: tomcat
      #container_version: 8 # 当前 tomcat 只支持 8 (默认)
```

### none / tomcat

```yaml
- java:
    params:
      build_type: none
      target: ${git-checkout}/user.war
      container_type: tomcat
```

如果需要复制文件到容器，可以使用以下方式

```yaml
# target也可以是url地址, 有额外的文件需要复制可以使用copy_assets
- java:
    params:
      build_type: none
      target: https://xxx.com/service.jar
      container_type: openjdk
      copy_assets:
        - xxx.html:/var/xxx.html
```

### gradle / spring-boot

```yaml
- java:
    params:
      jdk_version: 11
      build_type: gradle
      build_cmd: ./gradlew :user:build
      workdir: ${git-checkout}
      target: ./user/target/user.jar
      container_type: openjdk
      #container_version: 8 # 默认和 jdk_version 一致。可选: 8 / 11 / 17 / 21
```

## 自定义 JVM 参数

在 **环境部署 → 参数设置** 中，您可以通过添加 `JAVA_OPTS` 来自定义并追加 JVM 配置参数。  

如果希望完全替换系统默认的 `JAVA_OPTS`，请同时添加环境变量 `DISABLE_PRESET_JAVA_OPTS=true`。