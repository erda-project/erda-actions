# MYSQLDUMP

## 示例

Dice 某企业只有一个业务集群，该边缘集群切换中心集群。

```yaml
version: "1.1"
stages:
  - stage:
      - mysqldump:
          params:
            mysql_host: ((mysql_host))
            mysql_port: 3306
            mysql_username: ((mysql_username))
            mysql_password: ((mysql_password))
            mysql_database: dice
            dump_all_tables: "true"
            dump_all_tables_ignore_regexp: ^(fdp|pmp|okr|pipeline|dice_repo_web|dice_notify_histories|dice_nexus_|qa_sonar|qa_test_records|s_|ps_activities|ps_runtime_instances|ci_v3_build|ps_tickets|cm_containers|uc_user_event_log).*$
            global_drop_table_if_not_exists: true
            global_charset: utf8mb4
            global_try_filter_by_columns_and_values:
              org_id: ((org_id))
              cluster_name: ((cluster_name))
            must_include_tables:
              - table: ps_orgs
                where: id=((org_id))
              - table: admin_members
                where: scope_type='sys' OR org_id=((org_id))
              - table: dice_notify_groups
                where: org_id=((org_id))
            post_commands:
              - echo hello world
              - head ${DUMP_FILE_PATH}
```
