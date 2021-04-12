### Sonar Action

使用 sonar 进行代码质量检测

#### 使用

```yml
- stage:
  - sonar:
      params:
        code_dir: ${git-checkout}
```
