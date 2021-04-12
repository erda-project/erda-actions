## Android Action

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
      - android:
          params:
            context: ${repo}
            commands:
              - npm ci
              - cd android && ./gradlew clean assembleRelease
            target: build
```

