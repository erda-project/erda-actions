# Contributing

当你准备开发一款新的扩展或者对已存在的扩展进行维护升级时，请详细阅读此文档。

## 如何开发 Action

### 标准目录结构

一款 Action 的多个版本之间使用如下方式组织目录结构：

在每个具体版本目录下，除了包含基本的三大文件之外，还需要提供：

- internal 目录用于存放代码
- Dockerfile 文件用于制作镜像
- 在 Makefile 中添加相应的 TARGET

```bash
$ tree actions/echo
actions/echo
├── 0.1            # 0.1 版本
│   ├── internal   # 0.1 版本代码
│   │   └── .keep
│   ├── Dockerfile # Dockerfile 用于构建 0.1 版本的镜像
│   ├── README.md
│   ├── dice.yml
│   └── spec.yml
└── 1.0            # 1.0 版本
    ├── internal   # 1.0 版本代码
    │   └── .keep
    ├── Dockerfile # Dockerfile 用于构建 1.0 版本的镜像
    ├── README.md
    ├── dice.yml
    └── spec.yml
```

## 如何开发 AddOn

TODO

## 如何做好版本兼容

请查看 [这里](https://dice.app.terminus.io/workBench/projects/70/apps/178/repo/tree/develop/docs/dice-version-compatible.md)
