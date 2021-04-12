## redis-cli

根据选中的数据源执行 redis 命令

用法:

```yaml
- redis-cli:
    alias: redis-cli
    version: "1.0"
    params:
      command: "ping"
      datasource: "数据源的 id"
```