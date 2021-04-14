### Dice Version Archive
将 dice/dice 仓库下的 released dice.yml 和 migrations 打包迁移到 dice/version 仓库下。

- dice/dice:dice.yml --> dice/version:{version}/releases/dice/dice.yml
- dice/dice:.dice/migration --> dice/version:{version}/sqls

目录结构如下：

> dice repo
```text
{some_branch:repo_root}
├── .dice
│  └── migrations
│      ├── apim
│      │    ├── 210101_base.sql
│      │    └── 210201_some_feature.sql 
│      │    
│      └── cmdb
│           ├── 210101_base.sql
│           └── 210201_some_feature.sql
│
└── dice.yml
```

> version repo
```text
repo_root:version
   ├── {version_dir}
   ├── 3.21
   ├── 4.0
   │    ├── migrations
   │    │   └── srv
   │    │       └── v001
   │    │         ├── migrations.sql
   │    │         └── run.sh
   │    └── releases
   │    
   └── 4.1
        ├── migrations
        ├── releases
        │       └── dice
        │           └── dice.yml    <-- dice/dice:dice.yml 拷贝粘贴到这里
        └── sqls                    <-- dice/dice:.dice/migration 拷贝粘贴到这里 
              ├── apim
              │    ├── 210101_base.sql
              │    └── 210201_some_feature.sql
              │
              └── cmdb
                   ├── 210101_base.sql
                   └── 210201_some_feature.sql
```

