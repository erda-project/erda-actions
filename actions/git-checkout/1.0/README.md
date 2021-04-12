# Git-Checkout Action

## 简介

克隆 Git 代码仓库

## 详细介绍

git-checkout action 提供对 Git 代码仓库的克隆能力，包括设置 depth、submodules 等功能。

## params

git-checkout action 支持的配置。

### uri

代码仓库地址。

> 如果是在 Dice 界面上应用的流水线界面新建流水线，则平台会自动注入托管在平台上的该应用的代码仓库地址。

### branch

代码需要检出的分支。同时支持 commit / tag。

> 如果是在 Dice 界面上应用的流水线界面新建流水线，则平台会自动注入新建流水线时你选择的分支。

### username

若为私有仓库，需要提供鉴权信息。

> 如果是在 Dice 界面上应用的流水线界面新建流水线，则无需填写该字段。

### password

若为私有仓库，需要提供鉴权信息。

> 如果是在 Dice 界面上应用的流水线界面新建流水线，则无需填写该字段。

### depth

指定 `git clone --depth` 参数。

选填。

### submodules

是否获取 submodules。

选填。

可选值：

- none 不获取
- 具体子模块列表，例如：["moduleA","moduleB"]
- all 获取全部

### submodule_recursive

是否递归拉取 submodule。默认为 true。

选填。

可选值：

- true
- false

## outputs

outputs 可以通过 `${alias:OUTPUT:output}` 的方式被后续 action 引用。

支持的 outputs 列表如下：

- commit
- author
- author_date
- committer
- committer_date
- branch
- message

## 例子

```yml
- git-checkout:
    params:
      uri: https://github.com/xxx/yyy.git
      depth: 1

- custom-script:
    commands:
    - echo commit ${git-checkout:OUTPUT:commit} # 引用 git-checkout 产出的 output
```
