### PHP Action

#### 参数说明
 - index_path： php入口目录，一般是index.php的目录, 默认为context根目录

 - context: php代码目录

 - service: 在dice.yml中对应的服务名称
 
#### 使用

```yml
version: "1.1"
stages:
- stage:
  - git-checkout:
      alias: repo
      params:
        depth: 1
- stage:
  - php:
      params:
        index_path: web
        context: ${repo} 
        service: test-php
- stage:
  - release:
      params:
        dice_yml: ${repo}/dice.yml
        replacement_images:
        - ${php}/pack-result

- stage:
  - dice:
      params:
        release_id_path: ${release}
```