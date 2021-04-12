### Dockerfile Action

用于自定义 dockerfile 进行打包，产出应用镜像用于部署。

#### 使用

dockerfile 位于应用根目录:

```yml
- dockerfile:
  params:
    workdir: ${git-checkout}
    path: Dockerfile
    build_args:
      JAVA_OPTS: -Xms700m
      NODE_OPTIONS: --max_old_space_size=3040
```

dockerfile 位于应用子目录下:

```yml
- dockerfile:
  params:
    workdir: ${git-checkout}
    path: subdir/Dockerfile
```
