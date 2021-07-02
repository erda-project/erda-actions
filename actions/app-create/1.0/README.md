### 应用创建

创建内置仓库，并推送代码，当应用已经存在的时候将不会创建和推送

```yaml
version: "1.1"
stages:
  - stage:
      - app-create:
          alias: app-create
          version: "1.0"
          params:
            application_git_password: xxxx     # 仓库的密码
            application_git_repo: http://xxxx  # 仓库地址
            application_git_username: xxx      # 仓库的账号
            application_name: helloworld11 # 创建应用的名称
            application_type: SERVICE  # 创建应用的类型  SERVICE, SERVICE, MOBILE
            is_external_repo: false   
```

#### is_external_repo 说明

是否是外置仓库，如果是外置仓库，将采用上面填写的仓库信息来进行配置外置仓库引用信息

如果是内置仓库，则会根据上面配置的仓库信息拉取代码，然后推送到内置仓库中

