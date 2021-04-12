### Api Register Action

用于注册Api到网关

#### 使用

```yml
- api-register:
    params:
      release_id: ${release:OUTPUT:releaseID}
      swagger_json: ${swagger:OUTPUT:swaggerJson}
```
