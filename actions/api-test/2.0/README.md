### API-Test Action

执行单个接口测试。上层可以通过 pipeline.yml 编排一组接口测试的执行顺序。

#### 使用

```yaml
api-test:
  version: 2.0
  params:
    name: 访问 Dice 官网
    url: https://dice.terminus.io
    method: GET
    params:
      - key: p1
        value: v1
      - key: p2
        value: v2
    headers:
      - key: h1
        value: v1
      - key: h2
        value: v2
    body:
      type: application/json
      content: '{"name":"dice"}'
```