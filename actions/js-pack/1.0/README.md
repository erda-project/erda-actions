### js-pack Action

例子:

```yaml
- js-pack:
    alias: front-build
    params:
      node_version: 12 # node 版本，可选值 12、14，默认为 14
      build_cmd:
        - npm config set registry=https://registry.npm.terminus.io/ && npm i
        - npm run build
      workdir: ${{ dirs.git-checkout }}
      preserve_time: 600 # 报错时容器保留时间
```

在报错后，容器会按 `preserve_time` 配置时长继续运行。在日志顶部，会打印出类似 `NAMESPACE: pipeline-102679155835278` 这样的内容，复制好。

进入 **多云管理平台 > 容器资源 > Pods**，顶部选择分支对应的集群，然后在下方的命名空间里粘贴，应该可以查到该流水线。

点击该流水线记录后可以看到该流水线的容器，通过操作可进入容器控制台，然后就可以进行调试了。
