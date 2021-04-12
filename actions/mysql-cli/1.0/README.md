## mysql-cli

根据选中的数据源执行 sql 命令

用法:

```yaml
- mysql-cli:
    alias: mysql-cli
    version: "1.0"
    params:
      database: "database name"
      datasource: "数据源的 id"
      sql: "show tables"
```