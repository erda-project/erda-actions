### 执行应用流水线

根据应用名称分支和 yml 的名称来执行对应的流水线，然后等待执行完成

```yaml
version: "1.1"
stages:
  - stage:
      - app-run:
          alias: app-run
          version: "1.0"
          params:
            application_name: helloworld # 应用名称
            branch: master # 应用分支
            pipeline_yml_name: pipeline.yml  # 分支下流水线名称 例如: pipeline.yml  xxx.yml
```
