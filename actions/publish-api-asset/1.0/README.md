### Publish Api Asset Action

用于将 API 描述文档发布到 API 集市，成为 API 资源

#### 使用

example

```yml
- publish-api-asset:
    params:
      runtime_id: ${dice:OUTPUT:runtimeID}
      service_name: trade-center
      display_name: 交易接口
      asset_id: trade
      spec_path: ${git-checkout}/swagger.json
```
