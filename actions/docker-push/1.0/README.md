### docker-push Action

docker-push 主要完成如下两个功能：
1.从本地文件读取镜像，push 到指定镜像仓库
2.将 外部仓库的镜像 push 到 dice 内部仓库，供部署使用

#### 使用

Examples:

1. 从本地文件读取镜像，push 到指定镜像仓库
本地文件格式为：
module_name: test-server
image: xxxx

```yml
- docker-push:
    params:
      image: registry.erda.cloud/erdaxxx:v1.0   // 要 push 到外部镜像名称, require
      from: imageResult.img                               // 应用下面的文件
      service: test-server                                // 服务名称，要与镜像文件里的module_name一致
      username: admin                                     // 外部镜像仓库用户名
      password: xxxx                                      // 外部镜像仓库用密码
```

2. 将 外部仓库的镜像 push 到 dice 内部仓库，供部署使用

```yml
- docker-push:
    params:
      image: registry.erda.cloud/erdaxxx:v1.0   // 要 pull 的外部镜像名称, require
      service: test-server                                // 服务名称
      username: admin                                     // 外部镜像仓库用户名
      password: xxxx                                      // 外部镜像仓库用密码
      pull: true                                          // 该值必须为: true, require
```