### erda-get-service-addr Action

查看指定 runtimes 的 services 地址。

#### 使用

```yml
- erda-get-service-addr:
    params:
      runtime_id: ${runtime}/runtime-id
```
执行结果是输出一组 Meta 表示获取的 service 地址信息，对应的 Meta 的名称是服务名称，value 是服务的地址

后续 Action 可以通过 ${{ outputs.alias.val }} 获取并使用 Meta 信息。

例如要获取 service abc 的地址，则后续 Action 可以通过 ${{ outputs.abc.val }} 获取并使用 Meta 信息。