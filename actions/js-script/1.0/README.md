### JS-Script Action

提供js环境使用自定义命令构建通用产物

#### 使用


```yml
- js:
  params:
    workdir: ${git-checkout}
    commands: 
      - npm ci
      - npm build
      - npm run wechat:build
    targets:
      - out/web.zip
      - out/wechat.zip
```
