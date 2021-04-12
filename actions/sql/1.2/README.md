### SQL Action

### 一、简介

SQL Action执行用户配置的SQL脚本。
### 二、详细介绍

SQL Action负责将用户配置的SQL脚本提交至计算引擎执行，目前支持的计算引擎包括：Hive和Spark SQL。在执行SQL脚本时，SQL Action会检测相关依赖，保证工作流执行的有序性。
### 三、使用场景
SQL的使用场景是执行SQL语句进行ETL任务。目前支持的计算引擎是Hive和Spark SQL。

### 四、使用方式
```aidl
version: '1.0'
triggers:
  - schedule:
      cron: "0 */10 * * * ?"
resources:
- name: repo
  type: git
  source:
    uri: ((gittar.repo))
    branch: ((gittar.branch))
    username: ((gittar.username))
    password: ((gittar.password))
- name: sql-process
  type: sql
  source:
    queryType: sparksql
    queryEndPoint: jdbc:hive2://1.1.1.1:9000/;auth=noSasl
    username: foo
    password: foo
stages:
- name: repo
  tasks:
  - get: repo
    params:
      depth: 3
- name: sql-process
  tasks:
  - put: sql-process
    params:
      path: repo/ea/etl_terminus/dwd/rd/process/bug_df.q
      queryargs:
         PT_DATE: ${horus:getDateFromNow('','','yyyy-MM-dd',-1,'D')}
      inputTables:
      - pmp.s_pmp_issues
      - pmp.s_pmp_org_user_relatives
      - pmp.s_pmp_organizations
      - pmp.s_pmp_users
      - pmp.s_pmp_org_user_relatives
      - pmp.s_pmp_organizations
      - pmp.s_pmp_users
      - pmp.s_pmp_projects
      - pmp.s_pmp_workload_issue_statistics
      - pmp.s_pmp_issue_status
      outputTables:
      - ea.dwd_ea_bug_df
      triggerType: DAY
      frequency: 30
      process: process
```
