## iOS Action

#### 示例

```yml
version: "1.1"
stages:
  - stage:
      - git-checkout:
          alias: repo
          params:
            depth: 1
  - stage:
      - ios:
          params:
            context: ${repo}
            commands:
              - mkdir result
              - mv ${repo}/test.ipa result/test.ipa
            targets: 
              - result
  - stage:
      - release:
          params:
            release_mobile_file: 
                path: ${ios}/result/test.ipa
```

