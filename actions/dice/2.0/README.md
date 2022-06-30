# Dice Action

## 简介

提供部署服务的能力。

## params

### release_id_path

必填。

release 文件所在的目录，一般由前置的 release action 生成。

dice action 读取该 release 文件获取 releaseID，使用该 releaseID 调用 dice 平台开始应用部署。

例子：

```yaml
- dice:
    params:
      release_id_path: ${release}
```

### time_out

选填。

应用部署超时时间。单位为秒。默认 600 秒。

超过超时时间仍未部署成功，dice action 会主动调用取消部署接口，取消本次部署。

### deploy_without_branch

选填.

无分支名称部署, 项目级/应用级制品部署默认无分支名称，通过 release_id 部署时候可以选择该参数 (可选值：true/false)。

案例：

erda-demo 应用构建 develop 分支后，部署 runtime 名称为 `erda-demo/develop`, 使用该参数后，将部署/更新
分支对应环境下名称为 `erda-demo` 的 runtime。

## outputs

outputs 可以通过 ${alias:OUTPUT:output} 的方式被后续 action 引用。

支持的 outputs 列表如下：

- runtimeID

示例：

```yaml
- stage:
  - dice:
      ......
- stage:
  - custom-script:
      commands:
      - echo runtimeID: ${dice:OUTPUT:runtimeID}
```

## 例子

```yaml
- dice:
    params:
      release_id: ${release:OUTPUT:releaseID}
```
