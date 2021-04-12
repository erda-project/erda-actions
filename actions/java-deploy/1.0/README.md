### Java Library Publish Action

用于将用户的 java library 推送至远程 registry, 供其他用户使用

#### 使用

Examples:

- deploy 应用下所有模块

```yaml
- java-deploy:
    params:
      workdir: ${git-checkout}
      registry: ((REGISTRY))
      username: ((USERNAME))
      password: ((PASSWORD))
```

- deploy 应用下所有模块，并执行测试代码

```yaml
- java-deploy:
    params:
      workdir: ${git-checkout}
      registry: ((REGISTRY))
      username: ((USERNAME))
      password: ((PASSWORD))
      skip_tests: false
```

- deploy 应用下指定模块

```yaml
- java-deploy:
    params:
      workdir: ${git-checkout}
      registry: ((REGISTRY))
      username: ((USERNAME))
      password: ((PASSWORD))
      modules: moduleA,moduleB
```

- deploy 应用下指定模块，并执行测试代码

```yaml
- java-deploy:
    params:
      workdir: ${git-checkout}
      registry: ((REGISTRY))
      username: ((USERNAME))
      password: ((PASSWORD))
      modules: moduleA,moduleB
      skip_tests: false
```

- 用户指定 deploy 命令，例如使用 gradle 进行包发布

```yaml
- java-deploy:
    params:
      workdir: ${git-checkout}
      registry: ((REGISTRY))
      username: ((USERNAME))
      password: ((PASSWORD))
      modules: moduleA,moduleB
      skip_tests: false
      cmd: gradlew publish # 用户自己实现 gradle 脚本进行 publish
```
