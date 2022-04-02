# Multi-Merge Action

## 简介

在`git-checkout`action的基础上支持多个代码仓库的分支合并

## 详细介绍

multi-merge action 提供对 Git 代码仓库的克隆能力，并且支持多个代码仓库的分支合并。典型场景是我们提交pull request分支还没有合进去，但是我们需要合并这些分支到环境进行测试。

## params

multi-merge action 支持的配置。

### dest_repo

目标代码仓库地址。

> 如果是在 Erda 流水线界面新建流水线，则平台会自动注入托管在平台上的该应用的代码仓库地址。

### dest_branch

代码需要检出的分支。

> 如果是在 Erda 流水线界面新建流水线，则平台会自动注入新建流水线时你选择的分支。

### repos
所有需要merge的代码仓库分支列表。

### username

若为私有仓库，需要提供鉴权信息。

> 如果是在 Erda 流水线界面新建流水线，则无需填写该字段。

### password

若为私有仓库，需要提供鉴权信息。

> 如果是在 Erda 流水线界面新建流水线，则无需填写该字段。


## 例子

```yaml
version: "1.1"
stages:
  - stage:
      - multi-merge:
          alias: multi-merge
          description: 多代码仓库克隆
          version: "1.0"
          params:
            dest_branch: master
            dest_repo: https://github.com/xxx/yyy.git
            git_config:
              - name: http.https://github.com.proxy
                value: xxxx:8888
            repos:
              - branches:
                  - fix/xxx
                uri: https://github.com/aaa/aaa.git
              - branches:
                  - feat/xxx
                uri: https://github.com/bbb/bbb.git
  - stage:
      - custom-script:
          commands:
            - echo commit ${{ outputs.multi-merge.merged-repos }} # 引用 multi-merge 产出的 output
```
该例子会将repos下所有的分支都merge到master分支，并输出临时目录
