### JavaScript Action

用于编译打包 javascript 工程，产出应用镜像用于部署(Node.js版本: 12.13)。

#### outputs

outputs 可以通过 ${alias:OUTPUT:output} 的方式被后续 action 引用。

支持的 outputs 列表如下：

- image

示例：

```yaml
- stage:
  - js:
      ......
- stage:
  - release:
      params:
        image: 
          serviceNameInDiceYml: ${js:OUTPUT:image}
```

#### 使用

Examples:

1. spa 应用

```yml
- js:
  params:
    workdir: ${git-checkout}
    dependency_cmd: npm ci
    build_cmd: npm run build
    container_type: spa
    dest_dir: public
    npm_registry: <registry addr>
    npm_username: <npm username>
    npm_password: <npm password>
```

spa 应用用户须在应用根目录下放置 `nginx.conf.template` 文件, 文件内容如下:

```
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;

    # compression
    gzip on;
    gzip_min_length   2k;
    gzip_buffers      4 16k;
    gzip_http_version 1.0;
    gzip_comp_level   3;
    gzip_types        text/plain application/javascript application/x-javascript text/css application/xml text/javascript application/x-httpd-php image/jpeg image/gif image/png;
    gzip_vary on;

    client_max_body_size 0;

    location /api {
        proxy_pass http://<your-backend-addr>;
        proxy_set_header        X-Real-IP $remote_addr;
        proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header        Host $http_host;
    }
}
```

2. herd 应用

```yml
- js:
  params:
    workdir: ${git-checkout}
    build_cmd: npm run build
    container_type: herd
    dest_dir: public
```
