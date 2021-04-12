# MYSQLDUMP

## 示例

```yaml
version: "1.1"
stages:
  - stage:
      - mysqldump:
          params:
            # 数据库连接信息
            mysql_host: ((mysql_host))
            mysql_port: 3306
            mysql_username: ((mysql_username))
            mysql_password: ((mysql_password))
            mysql_database: db1
            # 导出所有表
            dump_all_tables: true
            # 在所有表基础上，过滤部分表，例如过滤 test1, test2 表
            dump_all_tables_ignore_regexp: ^(test1|test2).*$
            global_drop_table_if_not_exists: true
            global_charset: utf8mb4
            # 尝试使用字段和对应的值做过滤
            global_try_filter_by_columns_and_values:
              org_id: ((org_id))
              cluster_name: ((cluster_name))
            # 以下表一定会被导出
            must_include_tables:
              - table: ps_orgs
                where: id=((org_id))
              - table: admin_members
                where: scope_type='sys' OR org_id=((org_id))
              - table: dice_notify_groups
                where: org_id=((org_id))
            # 后置命令
            post_commands:
              - echo hello world
              - head ${DUMP_FILE_PATH}
```
