### Dockerfile Action

用于自定义 dockerfile 进行打包，产出应用镜像用于部署。

#### 使用

1. dockerfile 位于应用根目录:

```yml
- dockerfile:
  params:
    workdir: ${git-checkout}
    path: Dockerfile
    build_args:
      JAVA_OPTS: -Xms700m
      NODE_OPTIONS: --max_old_space_size=3040
```

2. dockerfile 位于应用子目录下:

```yml
- dockerfile:
  params:
    workdir: ${git-checkout}
    path: subdir/Dockerfile
```

3. 构建上下文配置

```yml
- stage:
    - git-checkout:
        alias: git-checkout
        version: "1.0"
        ...
    - git-checkout:
        alias: other-repo
        version: "1.0"
        ...
- stage:
    - dockerfile:
      params:
        build_context:
          other-resource: ${other-repo}
        workdir: ${git-checkout}
        path: subdir/Dockerfile
```

在 `Dockerfile` 中，使用 `COPY` 指令从不同的构建上下文（如 other-resource）中复制资源：

```Dockerfile
COPY --from=other-resource <source_path> <target_path>
```