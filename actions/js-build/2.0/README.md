### js-build Action

#### Params

| 属性         | 说明                     | 默认值 |
| ------------ | ------------------------ | ------ |
| node_version | 特殊指定运行的 node 版本 | 20     |
| build_cmd    | 需要执行的 sh 命令       | -      |
| workdir      | 指定执行命令的文件目录   | -      |

例子:

```yaml
- js-build:
    alias: js-build
    version: "2.0"
    params:
      node_version: 20
      build_cmd:
        - npm config set registry=https://registry.npm.terminus.io/ && npm i
        - npm run build
      workdir: ${git-checkout}
```

```yaml
- js-build:
    alias: js-build
    version: "2.0"
    params:
      build_cmd:
        - cnpm i
      workdir: ${git-checkout}
```

### release 用例

herd 模式

```yaml
- release:
    alias: release
    params:
      dice_yml: ${git-checkout}/dice.yml
      services:
        js-demo:
          cmd: cd /root/js-build && ls && npm run dev
          copys:
            - ${js-build}:/root/
          image: registry.erda.cloud/erda-actions/terminus-debian-herd:1.1.27-n20.15
```

spa 模式

```yaml
- release:
    alias: release
    params:
      dice_yml: ${git-checkout}/dice.yml
      services:
        js-demo:
          # 固定值，前提是项目中有 nginx.conf.template
          cmd: sed -i "s^server_name .*^^g" /etc/nginx/conf.d/nginx.conf.template && envsubst "`printf '$%s' $(bash -c "compgen -e")`" < /etc/nginx/conf.d/nginx.conf.template > /etc/nginx/conf.d/default.conf && /usr/local/openresty/bin/openresty -g 'daemon off;'
          # 固定值，注意 dist 是构建生成的产物
          copys:
            - ${js-build}/dist:/usr/share/nginx/html/ # dist 是使用 npm run build 生成出来的目录，常见的目录有：public、dist 等
            - ${js-build}/nginx.conf.template:/etc/nginx/conf.d/
          image: registry.erda.cloud/erda/terminus-nginx:0.2
```

