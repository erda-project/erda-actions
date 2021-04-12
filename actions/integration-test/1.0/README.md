# Integration-Test Action

提供针对 接口测试 的能力抽象

## 详细描述
it action 主要对应用进行接口测试，需要在 pipeline 里指定 it action 的相关信息进行接口测试，测试结果展示在应用详情的【应用测试】。

## params

### context

必填。
需要做it的代码存放目录。一般代码的编译为 git action 的 distination 目录。如repo。

### name

选填。
该次IT测试名称。

### parser_type

选填。
接口测试的框架类型，只支持 TESTNG/JUNIT。

#### 使用

```yml
- integration-test::
    params:
      code: ${git-checkout}/
```
