### Library Publish Action

用于发布库应用到应用市场，供用户引用

#### 使用

```yml
- lib-publish:
    params:
      workdir: ${git-checkout}
```

确保库应用源码根目录下有 `spec.yml` & `README.md` 文件
