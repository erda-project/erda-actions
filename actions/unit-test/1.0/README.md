### Unit Test Action

提供针对单元测试的能力抽象。

## 详细介绍
ut action 主要对用户的项目进行单元测试，当用户 push 代码时，会触发 ut-action，其中会探测应用的语言框架，选择相应的单测方式进行 unit test，测试结果展示在应用的【应用测试】。

## params

### context

必填。
需要做ut的代码存放目录。一般为 git action 的 destination 目录。如repo。若项目存在多种语言，必须指定模块路径，中间用 "," 分隔；如 "repo/path1,repo/path2"

### name

选填。
该次UT测试名称。

### context

选填。若UT的对象为golang，则必填。
该值为$GOPATH下的项目路径，如：terminus.io/dice/ci。

#### 使用

```yml
- unit-test:
    params:
      code: ${git-checkout}/
      command: ./gradlew test # 自定义单元测试命令，默认不用填写，平时自动分析语言类型并填充
```
