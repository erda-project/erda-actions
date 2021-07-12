# Erda Actions

## Extensions

Extension（扩展）是平台定义的扩展对象，用来自定义扩展平台的能力，目前包含两类扩展：Action 和 AddOn。

## 扩展的定义

扩展（以下简称 Ext ）使用 dice.yml、spec.yml 和 README.md 这三个文件来定义。这三个文件统称为一个 Ext 的 **元数据文件**。

### dice.yml

Ext 使用标准的 dice.yml 来描述 **运行时配置**。这里不再赘述 dice.yml 细节。

Action 和 AddOn 的 dice.yml 略有不同，请查看示例:

- [Action 示例 dice.yml](./actions/echo/1.0/dice.yml)
- [AddOn 示例 dice.yml](./addons/mysql/5.7.23/dice.yml)

### spec.yml

Ext 使用 spec.yml 来描述 **属性及参数规范**。

这里描述的信息主要用于 **界面图形化编辑** 以及 **官网扩展市场展示**。

```yaml
name: custom-script  # 名字
version: "1.0"       # 版本号
type: action         # 类型，是 action 还是 addOn
category: build      # 扩展所属的类别，用于在 Dice 官网扩展市场界面分类展示，目前可选：build / big_data / microservice / ability / search / distributed_cooperation / database / message
desc: execute cmds   # 对 Ext 的简单文字描述

# 参数声明，主要用于图形化编辑
# type: 参数类型，前端图形化编辑时会根据参数类型使用不同的输入框样式，并渲染出对应的文本。
	# 支持的类型列表，具体用法请参见下方示例
	# - string (默认值，不填写时即为 string 类型)
	# - string_array
	# - struct
	# - struct_array
	# - float
	# - int
	# - map
params:
  - name: command
    desc: 运行的命令
    type: string # string / float / int / map

  - name: commands
    desc: 一组命令
    type: string_array # go 代码中使用 []string 接收参数

  - name: person
    desc: 人
    type: struct
    struct:
      - name: name
        desc: 人名
        type: string # can omit
      - name: age
        desc: 年龄
        type: int

  - name: kvs
    desc: 键值对
    type: map # go 代码中使用 map[string]string 接收参数

accessibleAPIs: # 声明可访问的 OPENAPI 列表
  # push release
  - path: /api/releases
	method: POST
	schema: http

outputs: # 声明 action 的输出，可被后续 Action 通过 ${alias:OUTPUT:key} 方式引用
  - name: image
	desc: 输出的镜像
...
```

### README.md

README.md 是 Ext 的 **介绍文档**。

尽管你可以在 README.md 中填写任何内容，但是推荐的做法是写得越详细越好，因为这里的每一句话都能帮助用户更好地使用这款扩展。

## 内置扩展和第三方扩展

内置扩展由 Dice 平台开发者开发。与之对应的是第三方开发者开发的扩展，称为第三方扩展。

内置扩展和第三方扩展共同组成了扩展市场里的所有扩展。

## 扩展的使用

扩展市场里定义的扩展是可以运行的。Action 由 Pipeline 运行，AddOn 由 Orchestrator (AddOn-Platform) 运行。

### Action 如何使用

用户可以在 流水线描述文件 中引用 Action，同时在 Dice UI 代码仓库提供了流水线描述文件的图形化编辑界面，方便用户挑选和使用 Action。

### AddOn 如何使用

用户可以在 dice.yml 中引用 AddOn，同时在 Dice UI 代码仓库提供了 dice.yml 的图形化编辑界面，方便用户挑选和使用 AddOn。

## 内置扩展的管理

扩展市场内置了一批扩展，包括 Action 和 AddOn，分别在代码仓库 `actions/` 和 `addons/` 目录下。

### 内置扩展的版本如何管理

扩展市场的版本管理独立于 Dice 版本，不强跟随 Dice 版本升级而升级。

当前，扩展市场 **处于并将长期** 处于 1.0 版本，相对应的，该代码库只维护 master 分支，始终与 **端点 Dice 公有云平台** 的 [扩展市场](https://dice.terminus.io/market/pipeline) (以下简称 端点中心扩展市场) 保持强一致 (master 分支的任何改动都会被立刻同步至中心市场)。

这就要求每款扩展需要尽最大可能保持向前兼容。如果确实有无法兼容的改动出现，请先与团队沟通，再考虑是否升级这款扩展的版本号。

## 使用 dice cli 维护扩展

我们在 dice cli 中提供了 `dice ext` 子命令，方便对扩展进行维护。

### dice cli 的安装

dice cli 的安装，请查看 [这里](https://dice.app.terminus.io/workBench/projects/70/apps/178/repo/tree/develop/README.md)

### dice ext 命令在具体场景中的使用

dice ext 命令的简单介绍，请查看 [这里](#dice-ext)

#### 如何更改一个扩展使用的资源

- 通过 [dice ext pull](#dice-ext-pull) 拉取指定扩展的元数据文件
- 修改 dice.yml resources 字段
- 通过 [dice ext push](#dice-ext-push) 推送更新后的元数据文件至扩展市场，完成更新

#### 如何将市场 A 中的一款扩展推送至市场 B (无需推送镜像至 B 市场所在的 docker registry)

- 通过 `dice login` 登录到市场 A 所在云平台 (注意设置 `~/.dice.d/config` 的 server 配置)
- 通过 [dice ext pull](#dice-ext-pull) 拉取指定扩展的元数据文件
- 通过 `dice login` 登录到市场 B 所在云平台 (注意设置 `~/.dice.d/config` 的 server 配置)
- 通过 [dice ext push](#dice-ext-push) 将指定扩展推送至市场 B

#### 如何从一个市场 A 导出一个已经存在的扩展，并推送至另一个市场 B (同时推送镜像至 B 市场所在的 docker registry)

- 通过 `dice login` 登录到市场 A 所在云平台 (注意设置 `~/.dice.d/config` 的 server 配置)
- 通过 [dice ext export](#dice-ext-export) 导出指定扩展，包含 image blob
- 当前环境能否直接访问市场 B 及其 docker registry ?
  - 若可以，则跳过该步骤
  - 若不可以，则需要将元数据文件和 image blob 打包拷贝至一台可以直接访问市场 B 及其 docker registry 的机器上
- 通过 `dice login` 登录到市场 B 所在云平台 (注意设置 `~/.dice.d/config` 的 server 配置)
- 通过 [dice ext import](#dice-ext-import) 导入指定扩展，包括将 image blob 推送至私有 docker registry 上

## 安装扩展市场

该代码库始终与 **端点 Dice 公有云平台** 的 [扩展市场](https://dice.terminus.io/market/pipeline) 保持同步。

安装扩展市场时，所有扩展信息请直接通过 dice cli 登录 **端点 Dice 公有云平台** 后进行获取。

### 扩展市场在何时需要安装

扩展市场只在中心集群部署。因此，只有在添加中心集群时才需要安装扩展市场。添加 Edge 集群时无需关心扩展市场。

### Edge 集群与扩展市场的关系

Edge 集群不会部署扩展市场。

### 如何在全新安装的中心集群安装扩展市场

安装者根据所需要的扩展列表，将所需数据从 **端点中心扩展市场** 下载并安装至目标扩展市场。

每个 Dice 版本推荐的扩展市场列表，查看 [这里](https://dice.app.terminus.io/workBench/projects/70/apps/3287/repo/tree/develop/extensionctl)

历史 Dice 版本的扩展市场安装：

- Dice 3.4 的扩展市场安装，查看 [这里](docs/INSTALLATION-3.4.md)
- Dice 3.5 的扩展市场安装，查看 [这里](docs/INSTALLATION-3.5.md)

### 如何在中心集群升级扩展市场

目前不存在扩展市场升级的说法。

只需要将指定的扩展升级至指定版本即可。

## 如何开发一个扩展

查看 [这里](./docs/CONTRIBUTING.md)

## dice ext 命令简单介绍

#### dice ext

查询扩展

```bash
$ dice ext
ID                       TYPE     CATEGORY                  PUBLIC
api-gateway              addon    microservice              true
canal                    addon    database                  true
configcenter             addon    microservice              true
consul                   addon    distributed_cooperation   true
terminus-elasticsearch   addon    search                    true
kafka                    addon    message                   true
monitor                  addon    microservice              true
mysql                    addon    database                  true
redis                    addon    database                  true
registercenter           addon    microservice              true
...
```

#### dice ext pull

拉取指定 ext 的元信息

```bash
$ dice ext pull echo@1.0 -d echo
✔ extension pull success

$ ls echo
README.md dice.yml  spec.yml
```

#### dice ext push

> 需要系统管理员权限

推送指定 ext 的元信息

```bash
$ 对元信息进行更新

$ dice ext push -d echo --public -f
✔ extension echo push success
```

#### dice ext export

导出指定扩展，包含：

- 元信息
- 镜像 Blob

```bash
$ dice ext export echo@1.0 -d echo
registry.cn-hangzhou.aliyuncs.com login
 username: fb@terminus.io
password:
docker://registry.cn-hangzhou.aliyuncs.com/dice/echo-action => dir:echo/image/echo
Getting image source signatures
Copying blob sha256:e7c96db7181be991f19a9fb6975cdbbd73c65f4a2681348e63a141a2192a5f10
Copying blob sha256:6b34c3a5e37542504a101fcd7f324b0b36cdafeb13f1061e6ce2c58270201fac
Copying blob sha256:b9b8c261a200259323d418e12d7fb49950dd0fb11d5c18922ba692e3058c8ec5
Copying blob sha256:eb54ce0d56b748221c0c1b3353ef9a33cd4c8f36dd3542ba73b8254d4ffe9f61
Copying config sha256:373ef300e43acb288cb201c5999d2bc6f4b790d51652edc1bf9e93c8a56dc1a1
Writing manifest to image destination
Storing signatures
✔ extension echo@1.0 export success

$ tree -R echo
echo
├── README.md
├── dice.yml
├── image
│   └── echo
│       ├── 373ef300e43acb288cb201c5999d2bc6f4b790d51652edc1bf9e93c8a56dc1a1
│       ├── 6b34c3a5e37542504a101fcd7f324b0b36cdafeb13f1061e6ce2c58270201fac
│       ├── b9b8c261a200259323d418e12d7fb49950dd0fb11d5c18922ba692e3058c8ec5
│       ├── e7c96db7181be991f19a9fb6975cdbbd73c65f4a2681348e63a141a2192a5f10
│       ├── eb54ce0d56b748221c0c1b3353ef9a33cd4c8f36dd3542ba73b8254d4ffe9f61
│       ├── manifest.json
│       └── version
└── spec.yml

2 directories, 10 files
```

#### dice ext import

> 需要系统管理员权限

导入指定扩展，包括：

- 推送元信息至扩展市场
- docker push 镜像至 docker registry (若不指定 --registry 参数则不推送镜像)

```bash
$ dice ext import -d echo --registry localhost:5000
localhost:5000 login
 username:
password:
dir:echo/image/echo => docker://localhost:5000/dice/echo-action
Getting image source signatures
Copying blob sha256:e7c96db7181be991f19a9fb6975cdbbd73c65f4a2681348e63a141a2192a5f10
Copying blob sha256:6b34c3a5e37542504a101fcd7f324b0b36cdafeb13f1061e6ce2c58270201fac
Copying blob sha256:b9b8c261a200259323d418e12d7fb49950dd0fb11d5c18922ba692e3058c8ec5
Copying blob sha256:eb54ce0d56b748221c0c1b3353ef9a33cd4c8f36dd3542ba73b8254d4ffe9f61
Copying config sha256:373ef300e43acb288cb201c5999d2bc6f4b790d51652edc1bf9e93c8a56dc1a1
Writing manifest to image destination
Storing signatures
✔ extension echo@1.0 import success
```

## action list

### code
- [git-checkout](./actions/git-checkout)
- [git-push](./actions/git-push)

### build
- [gitbook](./actions/gitbook)
- [buildpack](./actions/buildpack)
- [maven-deploy](./actions/maven-deploy)
- [java](./actions/java)
- [java-build](./actions/java-build)
- [java-dependency-check](./actions/java-dependency-check)
- [java-deploy](./actions/java-deploy)
- [golang](./actions/golang)
- [js](./actions/js)
- [js-script](./actions/js-script)
- [js-build](./actions/js-build)
- [php](./actions/php)
- [dockerfile](./actions/dockerfile)
- [mobile-template](./actions/mobile-template)

### deploy

- [service-deploy](./actions/service-deploy)
- [dice](./actions/dice)
- [lib-publish](./actions/lib-publish)
- [mobile-publish](./actions/mobile-publish)
- [api-register](./actions/api-register)
- [publish-api-asset](./actions/publish-api-asset)
- [dice-deploy](./actions/dice-deploy)
- [dice-deploy-addon](./actions/dice-deploy-addon)
- [dice-deploy-domain](./actions/dice-deploy-domain)
- [dice-deploy-redeploy](./actions/dice-deploy-redeploy)
- [dice-deploy-release](./actions/dice-deploy-release)
- [dice-deploy-rollback](./actions/dice-deploy-rollback)
- [dice-deploy-service](./actions/dice-deploy-service)
- [erda-mysql-migration](./actions/erda-mysql-migration)
- [push-extensions](./actions/push-extensions)
- [archive-extensions](./actions/archive-extensions)
- [archive-release](./actions/archive-release)


### version
- [npm-publish](./actions/npm-publish)
- [release](./actions/release)

### test
- [integration-test](./actions/integration-test)
- [unit-test](./actions/unit-test)
- [api-test](./actions/api-test)
- [sonar](./actions/sonar)
- [mysqldump](./actions/mysqldump)
- [manual-review](./actions/manual-review)
- [mysql-cli](./actions/mysql-cli)
- [redis-cli](./actions/redis-cli)

### data governance
- [spark-upload](./actions/spark-upload)
- [sql](./actions/sql)
- [datax](./actions/datax)
- [sftp-access](./actions/sftp-access)
- [spark](./actions/spark)
- [k8sflink](./actions/k8sflink)
- [k8sspark](./actions/k8sspark)

### custom action
- [custom-script](./actions/custom-script)
- [oss-upload](./actions/oss-upload)
- [docker-push](./actions/docker-push)
- [loop](./actions/loop)
- [assert](./actions/assert)
- [jsonparse](./actions/jsonparse)

