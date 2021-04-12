### JavaScript Library Publish Action

用于将用户的 npm library 推送至远程 npm registry, 供其他用户使用

#### 使用

Examples:

```yml
- js-deploy:
  params:
    workdir: ${git-checkout}
    registry: ((publisher.npm.url))
    username: ((publisher.npm.username))
    password: ((publisher.npm.password))
```
其中 `username` & `password` 用于执行 `npm publish` 时使用，值为占位符，dice 平台会在运行时注入。
