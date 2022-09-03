# Buildpack Action

## 简介

buildpack action 完成代码工程的编译并制作镜像。

目前已支持[语言](#language)：

- Java
- Node.js
- Dockerfile

生成文件：

- [`build-result`](#build-result): 表示编译结果
- [`pack-result`](#pack-result): 表示镜像结果

## 详细介绍

该 action 是对打包所使用的 Buildpack 的抽象。

基本步骤：

- 识别代码语言，自动检测出合适的 bp
- 编译
- 利用编译结果制作镜像

buildpack-action 一次只能编译一个上下文目录(context)。如果一个应用有 n 个 context，则至少需要定义 n 个 buildpack-action。

## params

### context

必填。

表示执行构建的上下文目录。所有需要编译的 *模块* 必须在该目录下。

一般情况下 context 为平时执行 `mvn clean package` 或者 `npm ci` / `npm run build` 等基本编译命令的目录。

### modules

必填。

表示需要构建的模块列表，即微服务列表。

这里只需要填写最终需要作为微服务运行的模块列表。

例子：

```yaml
modules:
  # blog-web 对应 dice.yml 中 services 下的 blog-web 微服务。
  - name: blog-web
  # blog-service 对应 dice.yml 中 services 下的 blog-service 微服务。
  - name: blog-service
    path: blog-service/blog-service-impl
  # user-service 对应 dice.yml 中 services 下的 user-service 微服务。
  - name: user-service
    path: user-service/user-service-impl
```

子属性：

#### name

必填。

模块名。即 微服务 的名字。需要与 dice.yml 中 services 下的某一个 service 对应。

#### path

选填。

当 path 不填写时，path 默认为 name。

若应用为 java 应用，context + path 为 pom.xml 里定义的子模块的 pom.xml 所在路径。

若应用为单模块应用，context 已经包含模块信息，此时 path 需要指定为 `.` 。

#### image

选填。

```yaml
image:
  概述：modules的子属性，可用于定义自定义镜像名称，指定私有仓库账号；
  子属性：
    name: 必填
      概述：完整镜像名；
      细节：由三部分组成：`{ registry }/{ repository }:{ tag }`
      默认行为：如果不填，由平台生成的镜像名
    username: 选填
      概述：用于 docker login 指定 private registry 的用户名
      细节：当推送镜像到私有 docker registry 时，需要提供 username
    password: 选填
      概述：用户 docker login 指定 private registry 的密码
      细节：与 username 对应的 password
```

### language

选填。

指定需要使用的 buildpack 类型。

一般情况下，该字段无需用户手动填写。系统会根据 context 检测出合适的 bp。

当系统无法自动检测，或自动检测的结果不正确时，需要手动指定。

可选值：

- java
- node
- dockerfile

每个 language 对应一个或多个 [build_type](#buildtype) 和 [container_type](#containertype)。

### build_type

编译类型。

- java
  - maven
  - maven-edas
- node
  - npm
- dockerfile
  - dockerfile

### container_type

运行时类型。

- java
  - springboot
  - edas
- node
  - herd
  - spa
- dockerfile
  - dockerfile

### http_proxy

配置 http 代理

```yaml
buildpack:
  params:
    http_proxy: http://1.2.3.4:8888
    https_proxy: http://1.2.3.4:8888
```

### https_proxy

配置 https 代理

### bp_args

选填。

bp_args 指定构建过程中可以使用的 K/V 对，格式为 map[string]string。

用户可以在构建过程中通过环境变量的方式获取。

> key 和 value 都必须是 string 类型，必要时需要进行转换：true -> "true", 1 -> "1"

目前开放的参数：

1. MAVEN_OPTS

> build_type: maven / maven_edas
 
> 生效位置：`${MAVEN_OPTS} mvn clean pacakge ...`

2. MAVEN_EXTRA_ARGS

> build_type: maven / maven_edas

> 生效位置：`mvn clean package ${MAVEN_EXTRA_ARGS}`
 
一般用来指定 profile，例如 "-P dev"

3. NODE_OPTIONS

> build_type: npm

> 生效位置：`${NODE_OPTIONS} npm ci`
>
> 生效位置：`${NODE_OPTIONS} DICE_WORKSPACE=${DICE_WORKSPACE} npm run build`

4. DEP_CMD

> build_type: npm

默认值为 `npm ci`。

例如：`npm i` 或 `yum install -y gcc && npm ci`

5. WEBPACK_DLL_CONFIG

> build_type: npm

指定 webpack dll config 配置文件路径。

若为 true，则会执行 `DICE_WORKSPACE=${DICE_WORKSPACE} npm run dll`

6. PUBLIC_DIR

> container_type: spa

spa 应用在执行 `npm run build` 后生成的静态文件的目录。

默认值为 public。

7. only_build

选填。

是否只编译应用，不制作镜像。

默认为 false。

若为 true，则只会有 [`build-result`](#build-result) 文件，不会生成 [`pack-result`](#pack-result) 文件。

## Action 输出

以下文件会被放置在 Action 运行目录下：

### build-result

```json
[
  {
    "module_name": "galaxy-admin",
    "artifact_path": "bp-backend/app/galaxy-admin/app.jar"
  },
  {
    "module_name": "galaxy-web",
    "artifact_path": "bp-backend/app/galaxy-web/app.jar"
  }
]
```

### pack-result

```json
[
  {
    "module_name": "galaxy-admin",
    "image": "docker-registry.registry.marathon.mesos:5000/galaxy/galaxy-admin:v0.1"
  },
  {
    "module_name": "galaxy-web",
    "image": "docker-registry.registry.marathon.mesos:5000/galaxy/galaxy-web:v0.1"
  }
]
```

## 例子

### 一个应用同时包括前后端

```yaml
- stage:
  - buildpack:
      alias: bp-backend
      params:
        context: ${repo}/services/showcase
        modules:
        - name: blog-web
        - name: blog-service
          path: blog-service/blog-service-impl
        - name: user-service
          path: user-service/user-service-impl
  - buildpack:
      alias: bp-frontend
      params:
        context: ${repo}/endpoints/showcase-front
        modules:
        - name: showcase-front
```

### java (maven + springboot)

```yaml
- buildpack:
    params:
      language: java
      build_type: maven
      container_type: springboot
      context: ${repo} # your real context
      modules:
      - name: service1
        path: . # relative path according to context
```

### java (maven edas + edas)

edas 环境下，我们会改写你的 pom 文件，增加 edas-dubbo 依赖，并在 java 启动参数增加 edas 相关配置。

```yaml
- buildpack:
    params:
      language: java
      build_type: maven-edas
      container_type: edas
      context: ${repo} # your real context
      modules:
      - name: service1
        path: . # relative path according to context
```

### java 应用如何同时支持普通环境和 edas 环境

```yaml
- buildpack:
    params:
      language: java
      build_type: ((build_type)) # 私有配置，maven / maven-edas
      container_type: ((container_type)) # 私有配置，springboot / edas
      context: ${repo} # your real context
      modules:
      - name: service1
        path: . # relative path according to context
```

### node 应用 (herd)

```yaml
- buildpack:
    params:
      language: node
      build_type: npm
      container_type: herd
      context: ${repo} # your real context
      modules:
      - name: service1
        path: . # relative path according to context
```

### spa 单页应用

```yaml
- buildpack:
    params:
      language: node
      build_type: npm
      container_type: spa
      context: ${repo} # your real context
      modules:
      - name: service1
        path: . # relative path according to context
```

### 使用 Dockerfile 构建应用

```yaml
- buildpack:
    params:
      language: dockerfile
      context: ${repo} # your real context
      modules:
      - name: service1
        path: . # relative path according to context
```
