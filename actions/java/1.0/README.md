### Java Action

用于编译打包 java 工程，制作成为 docker image 用于运行服务。

#### outputs

outputs 可以通过 ${alias:OUTPUT:output} 的方式被后续 action 引用。

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

#### 使用

Examples:

1. maven / spring-boot

```yaml
- java:
    params:
      build_type: maven # 打包类型，这里选用 maven 打包
      workdir: ${git-checkout} # 打包时的当前路径，此路径一般为根 pom.xml 的路径
      options: -am -pl user # maven 打包参数，比如打包 user 模块使用命令 `mvn clean package -am -pl user`，这里省略命令 `mvn clean package` 只需要填写参数
      target: ./user/target/user.jar # 打包产物，一般为 jar，填写相较于 workdir 的相对路径。文件必须存在，否则将会出错。
      container_type: spring-boot # 运行 target（如 jar）所需的容器类型，比如这里我们打包的结果是 spring-boot 的 fat jar，故使用 spring-boot container
      #container_version: v1.8.0.181 # 可选: v1.8.0.181, v11.0.6, 默认 v1.8.0.181
```

```yaml
- stage:
  - java:
      alias: java-demo
      params:
        jdk_version: 11
        build_type: maven
```
2. maven / tomcat

```yaml
- java:
    params:
      jdk_version: 8
      build_type: maven
      workdir: ${git-checkout}
      options: -am -pl user
      target: ./user/target/user.war
      container_type: tomcat
      #container_version: v8.5.43-jdk8 # 可选: v8.5.43-jdk8, v7.0.96-jdk8, 默认 v8.5.43-jdk8
```

3. none / tomcat

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

4. gradle / spring-boot

```yaml
- java:
    params:
      jdk_version: 11
      build_type: gradle
      build_cmd: ./gradlew :user:build
      workdir: ${git-checkout}
      target: ./user/target/user.jar
      container_type: openjdk
      #container_version: v1.8.0.181 # 可选: v1.8.0.181, v11.0.6, 默认 v1.8.0.181
```