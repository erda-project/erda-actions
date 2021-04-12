### GitBook Action

用于打包gitbook

#### 使用

Examples:

pipeline.yml

```yaml
version: "1.1"
stages:
- stage:
  - git-checkout:
      params:
        depth: 1
- stage:
  - gitbook:
      params:
        workdir: ${git-checkout}
- stage:
  - release:
      params:
        dice_yml: ${git-checkout}/dice.yml
        image:
          gitbook: ${gitbook:OUTPUT:image}
- stage:
  - dice:
      params:
        release_id_path: ${release}
```

dice.yml

```yaml
version: 2
services:
  gitbook:
    ports:
      - 80
    resources:
      cpu: 1
      mem: 256
    expose:
      - 80
    hosts: []
    deployments:
      replicas: 1
    envs: {}
    binds: []
    health_check: {}
addons: {}
envs: {}

```