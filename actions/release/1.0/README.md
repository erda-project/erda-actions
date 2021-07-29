# Release Action

## 简介

release action 生成一个完整的 release (Dice 软件包) 用于发布。

## 详细描述

release action 根据 dice.yml 以及其他必要信息，生成一个完整的 release。

该 release 可以用来在 Dice 上进行部署服务。

用户提供的 dice.yml 中的 services 通常不会包括 image 字段，image 字段一般由 CI 构建后，通过 dicehub action 插入 dice.yml。

## params

### services

#### 例子

```yaml
 - release:
      params:
        dice_yml: ${git-checkout}/dice.yml
        services:
          java-demo: 
            image: "openjdk:8-jre-alpine"
            copys:
              - ${java-build:OUTPUT:buildPath}/target/docker-java-app-example.jar:/target
              - ${java-build:OUTPUT:buildPath}/spot-agent:/
            cmd: java ${java-build:OUTPUT:JAVA_OPTS} -jar /target/docker-java-app-example.jar
```

如例子所示，services 下需要配置和 dice.yaml 中的 service 一致，多个服务对应在 java-demo 同级下创建即可，
服务内部就有 cmd, image, cps 等参数需要填入，cmd 一般是你在本机环境执行的命令，而 image 对应就是执行的环境
，图形界面会优化，你只需要选择运行环境的版本即可，cps 就是将各个步骤的东西拷贝到执行环境中，不同语言的执行的例子可以到对应的 build
action 中去查看，比如 java-build 就有 spring-boot 的例子和 tomcat 的例子，上面的例子就是 spring-boot的例子

说明: 上面 spring-boot 的例子都会和本地有点不一样，就是需要将对应的东西拷贝到对应的位置，比如例子中的 jar 包就是要拷贝到 /target
中，然后还有些固定值 ${java-build:OUTPUT:JAVA_OPTS}，${java-build:OUTPUT:buildPath}/spot-agent/spot-agent.jar:/spot-agent/spot-agent.jar 
这些都是对应 build action 会说明的，一般加上这些固定值就可以接入平台的监控，当然不加也没事


### dice_yml

必填。

dice_yml 文件路径。

一般通过 ${git-checkout}/dice.yml 方式从代码仓库中进行引用。

例子：

```yaml
- release:
    params:
      dice_yml: ${git-checkout}/dice.yml
    ...
```

### replacement_images

选填。

当 dice.yml 需要用于服务部署时，需要保证每个 service 的 image 字段均不为空。

replacement_images 是一组文件列表，文件格式请参考 [链接](./buildpack.md#pack-result)。

例子：

```yaml
- release:
    params:
      replacement_images:
      - ${buildpack}/pack-result
      - ${frontend}/pack-result
    ...
```

### migration_dir

选填。

> * 用户可以指定一个包含数据库脚本的目录，目录下可以有多个 sql 文件，每个 sql 文件都严格依照flyway(一个数据库版本迁移工具)的标准。
> 
> * Flyway对数据库进行版本管理主要由Metadata表和6种命令完成，Metadata主要用于记录元数据，每种命令功能和解决的问题范围不一样，以下分别对metadata表和这些命令进行阐述，其中的示意图都来自Flyway的官方文档。
> 
> * Migrations是指Flyway在更新数据库时是使用的版本脚本，比如：一个基于Sql的Migration命名为V1__init_tables.sql，内容即是创建所有表的sql语句。

目录结构事例如下：
![](http://terminus-paas.oss-cn-hangzhou.aliyuncs.com/paas-doc/2020/03/03/e19af5bc-557d-4f09-8b8d-42c3ef596e8f.png)

对应的代码仓库下 db 目录的结构为：

```bash
$ tree ${repo}/db/migration
db
├── V1__custom.sql
├── V2__custom.sql
└── V3__custom.sql
```

Flyway对Migrations的扫描还必须遵从一定的命名模式，Migration主要分为两类：Versioned和Repeatable。

> Versioned migrations
>
> 一般常用的是Versioned类型，用于版本升级，每一个版本都有一个唯一的标识并且只能被应用一次，并且不能再修改已经加载过的Migrations，因为Metadata表会记录其Checksum值。其中的version标识版本号，由一个或多个数字构成，数字之间的分隔符可以采用点或下划线，在运行时下划线其实也是被替换成点了，每一部分的前导零会被自动忽略。
>
> Repeatable migrations
>
> Repeatable是指可重复加载的Migrations，其每一次的更新会影响Checksum值，然后都会被重新加载，并不用于版本升级。对于管理不稳定的数据库对象的更新时非常有用。Repeatable的Migrations总是在Versioned之后按顺序执行，但开发者必须自己维护脚本并且确保可以重复执行，通常会在sql语句中使用CREATE OR REPLACE来保证可重复执行。


> **默认情况下基于Sql的Migration文件的命令规则如下图所示：**
![](http://terminus-paas.oss-cn-hangzhou.aliyuncs.com/paas-doc/2020/03/03/e0663289-539d-4b27-b6bf-c94ebb56cfd1.png)

> 其中的文件名由以下部分组成，除了使用默认配置外，某些部分还可自定义规则。
> 
> * prefix: 可配置，前缀标识，默认值V表示Versioned，R表示Repeatable
> * version: 标识版本号，由一个或多个数字构成，数字之间的分隔符可用点.或下划线_
> * separator: 可配置，用于分隔版本标识与描述信息，默认为两个下划线__
> * description: 描述信息，文字之间可以用下划线或空格分隔
> * suffix: 可配置，后续标识，默认为.sql


**Flyway的其他更详细的标准可以参考[官网](https://flywaydb.org/documentation/migrations "Flyway说明").**

例子：

dice.yml
```yaml
mysql: #实例名称，使用共享实例需要自己填写:
  plan: mysql:basic
  options:
    create_dbs: flyway  #要创建的数据库，create if not exist，使用migration前必须要先创建好数据库
```

pipeline.yml
```yaml
- release:
    params:
      migration_dir: ${git-checkout}/db/migration   # 使用了migration功能，此选项必填，migration sql存放地址
      migration_mysql_database: flyway              # 使用了migration功能，此选项必填，用户指定需要migration的数据库
```


### dice_development_yml

选填。

dice_development.yml 文件路径。

当代码仓库中存在 dice_development.yml 时可以指定。

### dice_test_yml

选填。

同 [dice_development_yml](#dice_development_yml)

### dice_staging_yml

选填。

同 [dice_development_yml](#dice_development_yml)

### dice_production_yml

选填。

同 [dice_development_yml](#dice_development_yml)

### check_diceyml

选填。

是否需要对 dice.yml 进行格式校验。默认为 true。

在某些特定场景下，dice.yml 不面向发布，可能包含一些模板内容，无法通过校验，但仍然需要 release，则需要设置 `check_diceyml: false`。

### image

选填。

与 [replacement_images](#replacement_images) 功能相同，以字符串形式提供 微服务-镜像 对应信息。

例子：

```yaml
- release:
    params:
      image: |
      {"galaxy-admin":"image1","galaxy-web":"image2"}
```

### cross_cluster

选填。

生成的 release 是否可以跨集群部署。

```yaml
- release:
    params:
      cross_cluster: true # 生成的 release 可以跨集群部署。
```

## Action 输出

release action 运行成功后会在当前目录生成 `dicehub-release` 文件。

例子：

```text
86047fde0ad24eb1903498bd6ce58461

#### 使用

```yml
# 拉取代码
- stage:
    - git-checkout: # 没有特别的参数，可以省略 params，默认为当前代码仓库   

# 打包，构造 `镜像` (image)
- stage:
    - java:
        alias: bp-user # 多个 java-action 存在时，请使用 alias 区分
        params:
          build_type: maven # 打包类型，这里选用 maven 打包
          workdir: ${git-checkout} # 打包时的当前路径，此路径一般为根 pom.xml 的路径
          options: -am -pl user # maven 打包参数，比如打包 user 模块使用命令 `mvn clean package -am -pl user`，这里省略命令 `mvn clean package` 只需要填写参数
          target: ./user/target/user.jar # 打包产物，一般为 jar，填写相较于 workdir 的相对路径。文件必须存在，否则将会出错。
          container_type: spring-boot # 运行 target（如 jar）所需的容器类型，比如这里我们打包的结果是 spring-boot 的 fat jar，故使用 spring-boot container

    - java:
        alias: bp-item
        params:
          build_type: maven
          workdir: ${git-checkout}
          options: -am -pl item
          target: ./item/target/item.jar
          container_type: spring-boot
    
    - java:
        alias: bp-web
        params:
          build_type: maven
          workdir: ${git-checkout}
          options: -am -pl web
          target: ./web/target/web.jar
          container_type: spring-boot

# 构造 `版本` (release)
- stage:
    - release:
        params:
          dice_yml: ${git-checkout}/dice.yml # dice.yml 为部署内容的描述文件
          image: # 需要将打包产出的 image 填充进 dice.yml 中各个 services 中去，下面的 user, item, web 均为 dice.yml 中描述的各个 services
            user: ${bp-user:OUTPUT:image} # 选取此前 stage 中打包产出的镜像，使用 ${<name or alias>:OUTPUT:image} 的语法
            item: ${bp-item:OUTPUT:image}
            web: ${bp-web:OUTPUT:image}
```

#### 使用（针对老版本）

```yml
# 拉取代码
- stage:
    - git-checkout: # 没有特别的参数，可以省略 params，默认为当前代码仓库   

# 打包，构造 `镜像` (image)
- stage:
    - java:
        alias: bp-user # 多个 java-action 存在时，请使用 alias 区分
        params:
          build_type: maven # 打包类型，这里选用 maven 打包
          workdir: ${git-checkout} # 打包时的当前路径，此路径一般为根 pom.xml 的路径
          options: -am -pl user # maven 打包参数，比如打包 user 模块使用命令 `mvn clean package -am -pl user`，这里省略命令 `mvn clean package` 只需要填写参数
          target: ./user/target/user.jar # 打包产物，一般为 jar，填写相较于 workdir 的相对路径。文件必须存在，否则将会出错。
          container_type: spring-boot # 运行 target（如 jar）所需的容器类型，比如这里我们打包的结果是 spring-boot 的 fat jar，故使用 spring-boot container
          service: user # 这里需要和 dice.yml 一一匹配，打包构造镜像完成后，将此镜像填充进 dice.yml 的 services 中，这里的 user 为 dice.yml 中描述的其中一个 services

    - java:
        alias: bp-item
        params:
          build_type: maven
          workdir: ${git-checkout}
          options: -am -pl item
          target: ./item/target/item.jar
          container_type: spring-boot
    
    - java:
        alias: bp-web
        params:
          build_type: maven
          workdir: ${git-checkout}
          options: -am -pl web
          target: ./web/target/web.jar
          container_type: spring-boot

# 构造 `版本` (release)
- stage:
    - release:
        params:
          dice_yml: ${git-checkout}/dice.yml # dice.yml 为部署内容的描述文件
          replacement_images: # 需要将打包产出的 image 填充进 dice.yml 中各个 services 中去，下面的 user, item, web 均为 dice.yml 中描述的各个 services
            - ${bp-user}/pack-result # 使用 ${<name or alias>}/pack-result 的格式进行填充
            - ${bp-item}/pack-result
            - ${bp-web}/pack-result
```
