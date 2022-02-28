# Erda MySQL Migration

![](http://terminus-paas.oss-cn-hangzhou.aliyuncs.com/paas-doc/2022/02/18/9f85e163-3c4b-4685-ada5-d87eeafcdd80.png)

Erda MySQL 数据库迁移工具

## [方法论: 可持续集成的数据库版本控制](./Erda%20MySQL%20Migrator%3A%20Database%20Version%20Control%20for%20Continuous%20Integration.md)
[Erda MySQL Migrator: Database Version Control for Continuous Integration](./Erda%20MySQL%20Migrator%3A%20Database%20Version%20Control%20for%20Continuous%20Integration.md)

提要：
- 为什么要进行数据库表控制
- 如何使用 Erda MySQL Migration 工具进行数据库版本控制
- 如何使用 Python 脚本实现复杂数据迁移逻辑
- 其他特性：Dryrun，规约检查，文件篡改发现与修订机制，日志收集等

## 功能
该 Action 用于将代码仓库中的 SQLs 脚本更新到数据库中。
用户需要将用于 migration 的 SQLs 脚本提交到某个目录，并按模块进行分门别类。
如指定 `.erda/migrations` 目录为存放脚本的目录，那么目录结构为

```text
repo-root:.
├── .erda
│  └── migrations
│      ├── config.yml
│      ├── module_1
│      │    ├── 210101-base.sql
│      │    ├── 210101-feature-1.sql
│      │    └── 210201-feature-2.sql 
│      │    
│      └── module_2
│           ├── 210101_base.sql
│           └── 210201_some_feature.sql
├── other_directories
└── erda.yml
```
其中 module_1 和 module_2 是用户定义的业务模块名，可以自定义。
module 目录下存放 SQLs 脚本。

Erda MySQL Migration Action 会读取所有脚本，将安装到数据库中。

## 参数说明

| 参数                 | 说明                                                            |
|--------------------|---------------------------------------------------------------|
| workdir            | action 的工作目录，可以设置为 git-checkout 获取的仓库目录                       |             
| migrationdir       | migrations 物料所在的目录，action 通过拼接 workdir 和 migrationdir 定位到物料路径 | 
| database           | MySQL database 或 schema 名称                                    |
| skip_lint          | 是否跳过规约检查                                                      |  
| skip_sandbox       | 是否跳过沙盒预执行                                                     |  
| skip_pre_migration | 该参数已过期                                                        |  
| skip_migration     | 是否跳过执行 migration，为 true 时不执行任何 migration 操作                   |
| lint_config        | 规约检查配置文件                                                      |
| modules            | 指定要执行的 migration 的业务模块，如不指定，则执行所有模块                           |
| retry_timeout      | 连接数据库超时时间                                                     |
| mysql_host         | 第三方 MySQL 地址，不设置时默认为服务实例（runtime）所引用的 MySQL Addon 地址          |
| mysql_port         | 第三方 MySQL 端口，使用 MySQL Addon 时无须配置                             |
| mysql_username     | 第三方 MySQL 用户名，使用 MySQL Addon 时无需配置                            |
| mysql_password     | 第三方 MySQL 密码，使用 MySQL Addon 时无须配置                             |


### MySQL 设置
Action 可以从三处获取 MySQL 设置, 分别是 Action 参数, Pipeline 参数, 以及 Addon MySQL.
优先级是 Action > Pipeline > Addon MySQL.

#### 从 Action 参数获取 MySLQ 设置
在 stage 的 params 中设置 MySQL 参数, 拥有最高的优先级. 注意 `mysql_host`, `mysql_port`, `mysql_username`, `mysql_password`
中有任意一个参数是空值, 那么这一组参数都是无效的, action 会到下一优先级的参数中获取参数. 
```yaml
  - stage:
      - erda-mysql-migration:
          alias: erda-mysql-migration
          description: Erda MySQL Migration 工具
          version: "1.0"
          params:
            lint_config: .erda/migrations/config.yml
            migrationdir: .erda/migrations
            skip_pre_migration: true
            workdir: ${github}
            mysql_host: my.mysql.host         # mysql 参数
            mysql_port: 3306                  # mysql 参数
            mysql_username: root              # mysql 参数
            mysql_password: my-mysql-password # mysql 参数
            database: erda                    # mysql 参数
```

#### 从 Pipeline 参数获取 MySQL 设置
如果 action 没有在 Action 参数里获取到完整的 MySQL 设置, 则会尝试从流水线参数中获取 MySQL 设置.
相应的流水线参数名分别为:
- `migration_host`
- `migration_port`
- `migration_username`
- `migration_password`
- `migration_database`

#### 从 Addon MySQL 获取 MySQL 设置
如果 action 没有在 Action 参数和流水线参数中获取到完整的 MySQL 设置, 则会将对应环境下的 runtime
引用的 Addon MySQL 作为 MySQL 服务.
注意, 要使用的数据库名仍要在 Action 参数或流水线参数中配置.
如果 runtime 没有引用任何 Addon MySQL, 该 stage 会失败.

#### 注意事项
- 如果第一次使用该 Action 之前, 数据库中已经存在业务表了, 要将这部分业务表结构和初始化数据整理成基线 SQL 脚本并在脚本首行标记`# MIGRATION_BASE`。
- 指定的 `database` 如果不存在，该 Action 会自动创建。
- Action 执行模块内的 SQLs 脚本时，首先执行所有的标记了`# MIGRATION_BASE`的基线脚本，基线脚本有多个时，按字符序执行；然后执行其他脚本，其他脚本也是按字符序执行。所以为脚本命名时，务必注意按一定的字符序。建议命名方式为`日期+数字序号+feature描述`。脚本文件名后缀应当为`.sql`。
- Action 对脚本是增量执行的：每次执行前，都会比对执行记录，只执行上次执行后增量的脚本。执行过的脚本不应当修改内容或重命名，不然 Action 就比对不出哪些被执行过了。Action 的比对方式是在数据库中新增一个执行记录表`schema_migration_history`，将执行过的文件都记录下来，请不要删除改表。
- 该 Action 只允许执行 DDL(数据定义语言) 和 DML(数据操作语言)，不支持 TCL(事务控制语言) 和 DCL(数据控制语言)，所以该 Action 不允许脚本中存在事务控制、授权等操作。

## 规约配置

### 规约的基本结构
```yaml
- name: DefaultValueLinter
  alias: "默认值校验: created_at 默认值为 CURRENT_TIMESTAMP"
  white:
    patterns:
      - ".*-base$"
    committedAt:
      - <=20210101
    modules:
      - module_1
    filenames:
      - 20210202-my-feature.sql
  meta:
    columnName: created_at
    defaultValue: CURRENT_TIMESTAMP
```
`.name`

规约的名称，下文的配置文件示例中有所有的可用的规约的名称。

`.alias`

规约别名。当不配置时，以规约名称为别名。注意：规约别名是唯一的。

`.white`

白名单。白名单支持四种格式，当 SQL 脚本符合其中任意一项时，则跳过对该脚本的校验。

`.white.patterns` 

正则表达式列表，当 SQL 脚本名称符合列表中任意一条时，则跳过对该脚本的校验。

`.white.committedAt`

脚本名日期列表，当 SQL 脚本名称前 8 位为日期时，Action 将该日期与 `.white.committedAt` 中的所有条件进行比较，符合列表中任意一条时，则跳过对该脚本的校验。
如脚本名称为 "20201212-my-scirpt.sql" 会被以上配置跳过规约检查。
支持的运算符有 "=""<"">""<="">="。

`white.modules`

表示该模块下所有脚本都不作此规约校验。

`white.filenames`

该列表中的文件不作此规约校验。

`.meta`

规约的元信息，不同的规约有不同的元信息，下文会注意介绍。

### 规约：SQL 句型启用禁用清单
```yaml
- name: AllowedStmtLinter
  meta:
    - stmtType: AlterTableStmt.AlterTableOption
      forbidden: true
      white:
        filenames:
          - "20211115-hepa-domain-add-org-1"
```
其 meta 是一个 SQL 句型列表。其中 `stmtType` 表示句型，`forbidden` 表示是否禁用，当设置为 true 时，认为脚本中出现的该句型为非法句型。
`white` 表示白名单，符合白名单的脚本中出现的禁用句型不会被认为非法。`white` 的结构与基本结构中的 `white` 类似，但不支持 modules 列表。

### 规约：布尔类型字段名称与类型校验
```yaml
- name: BooleanFieldLinter
```
布尔类型的字段名称应当符合如"is_deleted""has_child"的结构。且布尔类型字段类型应当为"Boolean"或"Tinyint(1)"。

该规约无 meta 结构。

### 规约：字段名称校验
```yaml
- name: ColumnNameLinter
  meta:
    patterns:
      - "^[0-9a-z_]{1,64}$"
```
其 meta 给定一个正则表达式列表，字段名至少应当符合其中一个正则表达式。

### 规约：字段类型校验
```yaml
- name: ColumnTypeLinter
  meta:
    columnName: id
    types:
      - type: varchar
        flen: 36
      - type: char
        flen: 36
```
该规约规定了给定的字段的类型。如上述示例中规定了 id 字段的类型是 varchar(36) 或 char(36) 中的一种。

### 规约：完整 INSERT 语句校验
```yaml
- name: CompleteInsertLinter
```
该规约要求 INSERT 语句不可省略字段名。

### 规约：默认值校验
```yaml
- name: DefaultValueLinter
  meta:
    columnName: created_at
    defaultValue: CURRENT_TIMESTAMP
```
其 meta 指定一个字段名和默认值，表示在建表时，该字段应该为此默认值。

### 规约：小数校验
```yaml
- name: FloatDoubleLinter
```
不可以用 float 和 double 表示小数，应当用 decima。该规约没有 meta 结构。

### 规约：禁止使用外键
```yaml
- name: ForeignKeyLinter
```
该规约没有 meta 结构。

### 规约：索引长度校验
```yaml
- name: IndexLengthLinter
```
单列索引长度不得超过 767，联合索引长度不得超过 3072. 以上长度基于 utf8mb4 计算。

该规约没有 meta 结构。

### 规约：索引名称校验
```yaml
- name: IndexNameLinter
  meta:
    indexPattern: "^idx_.*"
    uniqPattern: "^uk_.*"
```
该规约规定了索引名称的格式。`indexPattern` 规定普通索引的名称格式，`uniqPattern` 规定了唯一索引的名称格式。

### 规约：MySQL 关键字作为表名字段名校验
```yaml
- name: KeywordsLinter
  meta:
    "ALL": true
```
其 meta 是一个 map 结构，表示该字符串为 MySQL 关键字或某种保留字，不可用于表名和字段名。

### 规约：禁止显示为时间字段插入值
```yaml
- name: ManualTimeSetterLinter
  meta:
    columnName: created_at
```
其 meta 设置一个字段名称，表示不得在 migration 为为该字段插入值。

### 规约：必要字段校验
```yaml
- name: NecessaryColumnLinter
  meta:
    columnName:
      - id
```
其 meta 规定了一系列字段名，表示在建表语句中，至少要有列表中的一个字段。

### 规约：必要列选项校验
```yaml
- name: NecessaryColumnOptionLinter
  meta:
    columnOptionType:
      - "ColumnOptionComment"
- name: NecessaryColumnOptionLinter
  alias: "updated_at 应当自动跟踪更新时间"
  meta:
    columnName: "updated_at"
    columnOptionType:
      - "ColumnOptionOnUpdate"
```
其 meta 规定了一个列选项列表，表示在建表语句中，所有的列至少有列表中的一个选项。
如果规定了 columnName，则表示规约仅对该列生效。
可配置的 columnOptionType 有：
```yaml
- "ColumnOptionNoOption"      
- "ColumnOptionPrimaryKey"    
- "ColumnOptionNotNull"       
- "ColumnOptionAutoIncrement" 
- "ColumnOptionDefaultValue"  
- "ColumnOptionUniqKey"       
- "ColumnOptionNull"          
- "ColumnOptionOnUpdate"      
- "ColumnOptionFulltext"      
- "ColumnOptionComment"       
- "ColumnOptionGenerated"     
- "ColumnOptionReference"     
- "ColumnOptionCollate"       
- "ColumnOptionCheck"         
- "ColumnOptionColumnFormat"  
- "ColumnOptionStorage"       
- "ColumnOptionAutoRandom"    
```

### 规约：必要建表选项校验
```yaml
- name: NecessaryTableOptionLinter
  alias: "必要的 tableOption: 表应当有 comment"
  meta:
    key: TableOptionComment
- name: NecessaryTableOptionLinter
  alias: "必要的 tableOption: 表应当注明 charset 为 utf8 或 utf8mb4"
  white:
    patterns:
      - ".*-base$"
  meta:
    key: TableOptionCharset
    values:
      - "utf8"
      - "utf8mb4"
```
其 meta 规定了建表语句必要的表选项。key 表示选项的名称，values 表示可选的值。有的表选项是没有值的。

支持的 tableOption key 有：
```yaml
- TableOptionNone
- TableOptionEngine
- TableOptionCharset
- TableOptionCollate
- TableOptionAutoIdCache
- TableOptionAutoIncrement
- TableOptionAutoRandomBase
- TableOptionComment
- TableOptionAvgRowLength
- TableOptionCheckSum
- TableOptionCompression
- TableOptionConnection
- TableOptionPassword
- TableOptionKeyBlockSize
- TableOptionMaxRows
- TableOptionMinRows
- TableOptionDelayKeyWrite
- TableOptionRowFormat
- TableOptionStatsPersistent
- TableOptionStatsAutoRecalc
- TableOptionShardRowID
- TableOptionPreSplitRegion
- TableOptionPackKeys
- TableOptionTablespace
- TableOptionNodegroup
- TableOptionDataDirectory
- TableOptionIndexDirectory
- TableOptionStorageMedia
- TableOptionStatsSamplePages
- TableOptionSecondaryEngine
- TableOptionSecondaryEngineNull
- TableOptionInsertMethod
- TableOptionTableCheckSum
- TableOptionUnion
- TableOptionEncryption
```

### 规约：主键字段名校验
```yaml
- name: PrimaryKeyLinter
  white:
    patterns:
      - ".*-base$"
  meta:
    columnName: id
```
其 meta 提供了一个 columnName 字段，表示该字段必须是主键，主键必须是该字段。

### 规约：表名校验
```yaml
- name: TableNameLinter
  meta:
    patterns:
      - "^erda_[a-z0-9_]{1,59}"
```
其 meta 提供了一个正则表达式列表，表示表名只要要符合其中一个表达式。

### 规约：varchar 长度校验
```yaml
- name: VarcharLengthLinter
```
varchar 长度不得超过 5000，超过 5000 请用 text 类型。该规约没有 meta 结构。

### 一个包含以上所有规约完整的实例

```yaml
# SQL 句型启用禁用清单
- name: AllowedStmtLinter
  white:
    patterns:
      - "^.*-base$"
  meta:
    - stmtType: CreateDatabaseStmt
      forbidden: true
    - stmtType: AlterDatabaseStmt
      forbidden: true
    - stmtType: DropDatabaseStmt
      forbidden: true
    - stmtType: DropUserStmt
      forbidden: true
    - stmtType: CreateTableStmt
    - stmtType: DropTableStmt
      forbidden: true
    - stmtType: DropSequenceStmt
      forbidden: true
    - stmtType: RenameTableStmt
      forbidden: true
    - stmtType: CreateViewStmt
      forbidden: true
    - stmtType: CreateSequenceStmt
      forbidden: true
    - stmtType: CreateIndexStmt
    - stmtType: DropIndexStmt
    - stmtType: LockTablesStmt
      forbidden: true
    - stmtType: UnlockTablesStmt
      forbidden: true
    - stmtType: CleanupTableLockStmt
      forbidden: true
    - stmtType: RepairTableStmt
      forbidden: true
    - stmtType: TruncateTableStmt
      forbidden: true
    - stmtType: RecoverTableStmt
      forbidden: true
    - stmtType: FlashBackTableStmt
      forbidden: true
    - stmtType: AlterTableStmt
    - stmtType: AlterTableStmt.AlterTableOption
      forbidden: true
      white:
        filenames:
          - "20211115-hepa-domain-add-org-1"
          - "20211222-dice-member-character"
          - "20211026-pipeline-cron-row-format"
          - "20210823-project-charset-utf8mb4"
    - stmtType: AlterTableStmt.AlterTableAddColumns
    - stmtType: AlterTableStmt.AlterTableAddConstraint
    - stmtType: AlterTableStmt.AlterTableDropColumn
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableDropPrimaryKey
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableDropIndex
    - stmtType: AlterTableStmt.AlterTableDropForeignKey
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableModifyColumn
    - stmtType: AlterTableStmt.AlterTableChangeColumn
      forbidden: true
      white:
        filenames:
          - "20220113-pipeline-definition-change-field"
          - "20220211-deployment-order-batches.sql"
    - stmtType: AlterTableStmt.AlterTableRenameColumn
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableRenameTable
      forbidden: true
      white:
        filenames:
          - "20220110-file-record-soft-delete"
          - "20211227-project-soft-delete"
    - stmtType: AlterTableStmt.AlterTableAlterColumn
    - stmtType: AlterTableStmt.AlterTableLock
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableAlgorithm
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableRenameIndex
    - stmtType: AlterTableStmt.AlterTableForce
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableAddPartitions
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableCoalescePartitions
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableDropPartition
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableTruncatePartition
      forbidden: true
    - stmtType: AlterTableStmt.AlterTablePartition
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableEnableKeys
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableDisableKeys
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableRemovePartitioning
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableWithValidation
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableWithoutValidation
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableSecondaryLoad
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableSecondaryUnload
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableRebuildPartition
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableReorganizePartition
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableCheckPartitions
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableExchangePartition
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableOptimizePartition
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableRepairPartition
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableImportPartitionTablespace
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableDiscardPartitionTablespace
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableAlterCheck
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableDropCheck
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableImportTablespace
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableDiscardTablespace
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableIndexInvisible
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableOrderByColumns
      forbidden: true
    - stmtType: AlterTableStmt.AlterTableSetTiFlashReplica
      forbidden: true
    - stmtType: SelectStmt
    - stmtType: UnionStmt
    - stmtType: LoadDataStmt
      forbidden: true
    - stmtType: InsertStmt
    - stmtType: DeleteStmt
    - stmtType: UpdateStmt
    - stmtType: ShowStmt
    - stmtType: SplitRegionStmt
      forbidden: true

# 布尔值命名与类型校验
- name: BooleanFieldLinter
  white:
    patterns:
      - ".*-base$"

# 必须设置表 character set
# 可以用 NecessaryTableOptionLinter 替代
- name: CharsetLinter
  white:
    patterns:
      - ".*-base$"
  meta:
    tableOptionCharset:
      utf8: true
      utf8mb4: true

# 列命名校验
- name: ColumnNameLinter
  white:
    patterns:
      - ".*-base$"
  meta:
    patterns:
      - "^[0-9a-z_]{1,64}$"

# id 字段类型校验
- name: ColumnTypeLinter
  alias: "字段类型校验: id 应当为 varchar(36) 或 char(36)"
  white:
    patterns:
      - ".*-base$"
    committedAt:
      - "<=20220215"
  meta:
    columnName: id
    types:
      - type: varchar
        flen: 36
      - type: char
        flen: 36

# created_at 类型校验
- name: ColumnTypeLinter
  alias: "字段类型校验: created_at 应当为 datetime 类型"
  white:
    patterns:
      - ".*-base$"
  meta:
    columnName: created_at
    types:
      - type: datetime

# updated_at 类型校验
- name: ColumnTypeLinter
  alias: "字段类型校验: updated_at 应当为 datetime 类型"
  white:
    patterns:
      - ".*-base$"
  meta:
    columnName: updated_at
    types:
      - type: datetime

# complete insert 校验
- name: CompleteInsertLinter
  white:
    patterns:
      - ".*-base$"

# created_at 字段默认值校验
- name: DefaultValueLinter
  alias: "默认值校验: created_at 默认值为 CURRENT_TIMESTAMP"
  white:
    patterns:
      - ".*-base$"
  meta:
    columnName: created_at
    defaultValue: CURRENT_TIMESTAMP

# updated 默认值校验
- name: DefaultValueLinter
  alias: "默认值校验: updated_at 默认值为 CURRENT_TIMESTAMP"
  white:
    patterns:
      - ".*-base$"
  meta:
    columnName: updated_at
    defaultValue: CURRENT_TIMESTAMP

# explicit_collation 校验
- name: ExplicitCollationLinter
  white:
    patterns:
      - ".*-base$"

# 小数类型校验
- name: FloatDoubleLinter
  white:
    patterns:
      - ".*-base$"

# 外键校验
- name: ForeignKeyLinter
  white:
    patterns:
      - ".*-base$"
    modules:
      - "fdp"
      - "fdp-agent"

# 索引长度校验
- name: IndexLengthLinter
  white:
    patterns:
      - ".*-base$"

# 索引名称校验
- name: IndexNameLinter
  white:
    patterns:
      - ".*-base$"
  meta:
    indexPattern: "^idx_.*"
    uniqPattern: "^uk_.*"

# MySQL 关键字 保留字检查: 不能将关键字作为表名或字段名
- name: KeywordsLinter
  white:
    patterns:
      - ".*-base$"
  meta:
    "ALL": true
    "ALTER": true
    "AND": true
    "ANY": true
    "AS": true
    "ENABLE": true
    "DISABLE": true
    "ASC": true
    "BETWEEN": true
    "BY": true
    "CASE": true
    "CAST": true
    "CHECK": true
    "CONSTRAINT": true
    "CREATE": true
    "DATABASE": true
    "DEFAULT": true
    "COLUMN": true
    "TABLESPACE": true
    "PROCEDURE": true
    "FUNCTION": true
    "DELETE": true
    "DESC": true
    "DISTINCT": true
    "DROP": true
    "ELSE": true
    "EXPLAIN": true
    "EXCEPT": true
    "END": true
    "ESCAPE": true
    "EXISTS": true
    "FOR": true
    "FOREIGN": true
    "FROM": true
    "FULL": true
    "GROUP": true
    "HAVING": true
    "IN": true
    "INDEX": true
    "INNER": true
    "INSERT": true
    "INTERSECT": true
    "INTERVAL": true
    "INTO": true
    "IS": true
    "JOIN": true
    "KEY": true
    "LEFT": true
    "LIKE": true
    "LOCK": true
    "MINUS": true
    "NOT": true
    "NULL": true
    "ON": true
    "OR": true
    "ORDER": true
    "OUTER": true
    "PRIMARY": true
    "REFERENCES": true
    "RIGHT": true
    "SCHEMA": true
    "SELECT": true
    "SET": true
    "SOME": true
    "TABLE": true
    "THEN": true
    "TRUNCATE": true
    "UNION": true
    "UNIQUE": true
    "UPDATE": true
    "VALUES": true
    "VIEW": true
    "SEQUENCE": true
    "TRIGGER": true
    "USER": true
    "WHEN": true
    "WHERE": true
    "XOR": true
    "OVER": true
    "TO": true
    "USE": true
    "REPLACE": true
    "COMMENT": true
    "COMPUTE": true
    "WITH": true
    "GRANT": true
    "REVOKE": true
    "WHILE": true
    "DO": true
    "DECLARE": true
    "LOOP": true
    "LEAVE": true
    "ITERATE": true
    "REPEAT": true
    "UNTIL": true
    "OPEN": true
    "CLOSE": true
    "CURSOR": true
    "FETCH": true
    "OUT": true
    "INOUT": true
    "LIMIT": true
    "DUAL": true
    "FALSE": true
    "IF": true
    "KILL": true
    "TRUE": true
    "BINARY": true
    "SHOW": true
    "CACHE": true
    "ANALYZE": true
    "OPTIMIZE": true
    "ROW": true
    "BEGIN": true
    "DIV": true
    "MERGE": true
    "PARTITION": true
    "CONTINUE": true
    "UNDO": true
    "SQLSTATE": true
    "CONDITION": true
    "MOD": true
    "CONTAINS": true
    "RLIKE": true
    "FULLTEXT": true

# 显式插入时间校验
- name: ManualTimeSetterLinter
  alias: "禁止显示插入 created_at 的值"
  white:
    patterns:
      - ".*-base$"
  meta:
    columnName: created_at

# 显式插入时间校验
- name: ManualTimeSetterLinter
  alias: "禁止显示插入 updated_at 的值"
  white:
    patterns:
      - ".*-base$"
  meta:
    columnName: updated_at

# 必要字段 id
- name: NecessaryColumnLinter
  alias: 必要字段.id
  white:
    patterns:
      - ".*-base$"
  meta:
    columnName:
      - id

# 必要字段 created_at
- name: NecessaryColumnLinter
  alias: 必要字段.created_at
  white:
    patterns:
      - ".*-base$"
  meta:
    columnName:
      - created_at

# 必要字段 updated_at
- name: NecessaryColumnLinter
  alias: 必要字段.updated_at
  white:
    patterns:
      - ".*-base$"
  meta:
    columnName:
      - updated_at

# 必要字段 org_name
- name: NecessaryColumnLinter
  alias: 必要字段.org_name
  white:
    patterns:
      - ".*-base$"
    committedAt:
      - "<=20220215"
  meta:
    columnName:
      - org_name

# 必要字段 org_id
- name: NecessaryColumnLinter
  alias: 必要字段.org_id
  white:
    patterns:
      - ".*-base$"
    committedAt:
      - "<=20220215"
  meta:
    columnName:
      - org_id

# 必要字段 deleted_at 或 soft_deleted_at
- name: NecessaryColumnLinter
  alias: 必要字段.deleted_at或soft_deleted_at
  white:
    patterns:
      - ".*-base$"
    committedAt:
      - "<=20220215"
  meta:
    columnName:
      - deleted_at
      - soft_deleted_at

# 必要列选项: comment
- name: NecessaryColumnOptionLinter
  alias: "列应当有注释"
  white:
    patterns:
      - ".*-base$"
  meta:
    columnOptionType:
      - "ColumnOptionComment"

# not null 校验: 所有字段应当显示定义为 not null
- name: NecessaryColumnOptionLinter
  alias: "列要么是主键要么是 not null 的"
  white:
    patterns:
      - ".*-base$"
  meta:
    columnOptionType:
      - "ColumnOptionNotNull"
      - "ColumnOptionPrimaryKey"

# updated_at 跟踪更新时间校验
- name: NecessaryColumnOptionLinter
  alias: "updated_at 应当自动跟踪更新时间"
  white:
    patterns:
      - ".*-base$"
  meta:
    columnName: "updated_at"
    columnOptionType:
      - "ColumnOptionOnUpdate"

# 必要的 tableOption：所有表应当有表注释
- name: NecessaryTableOptionLinter
  alias: "必要的 tableOption: 表应当有 comment"
  white:
    patterns:
      - ".*-base$"
  meta:
    key: TableOptionComment

# 必要的 tableOption: 表应当注明 charset
- name: NecessaryTableOptionLinter
  alias: "必要的 tableOption: 表应当注明 charset 为 utf8 或 utf8mb4"
  white:
    patterns:
      - ".*-base$"
  meta:
    key: TableOptionCharset
    values:
      - "utf8"
      - "utf8mb4"

# id 必须是主键, 主键必须是 id
- name: PrimaryKeyLinter
  white:
    patterns:
      - ".*-base$"
  meta:
    columnName: id

# 表名称校验 表名应当以 erda 开头
- name: TableNameLinter
  alias: "TableNameLinter: 以 erda_ 开头仅包含小写英文字母数字下划线"
  white:
    patterns:
      - ".*-base$"
    committedAt:
      - "<=20220215"
  meta:
    patterns:
      - "^erda_[a-z0-9_]{1,59}"

# varchar 长度校验
- name: VarcharLengthLinter
  white:
    patterns:
      - ".*-base$"
```