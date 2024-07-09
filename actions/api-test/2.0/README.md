### API-Test Action

执行单个接口测试。上层可以通过 pipeline.yml 编排一组接口测试的执行顺序。

#### 使用

推荐通过图形化界面进行编辑。

通过声明出参，可以将接口返回的数据传递给下游 action，例如下例中的 `a-cookie`。

```yaml
version: "1.1"
stages:
  - stage:
      - api-test:
          description: 执行单个接口测试。上层可以通过 pipeline.yml 编排一组接口测试的执行顺序。
          version: "2.0"
          params:
            asserts:
              - arg: a-cookie
                operator: not_empty
            body:
              type: none
            method: GET
            out_params:
              - expression: Set-Cookie
                key: a-cookie
                source: header
            url: https://www.erda.cloud/api/users/me
          timeout: 3600
  - stage:
      - custom-script:
          alias: cookie-printer
          description: 运行自定义命令
          version: "2.0"
          commands:
            - |-
              c="${{ outputs.api-test.a-cookie }}"
              echo "I got a cookie: $c"
            - 'echo "action meta: b-cookie=$c"'
          resources:
            cpu: 0.1
            mem: 256
```
