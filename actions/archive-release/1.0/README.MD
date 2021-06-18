<!--
 Copyright (c) 2021 Terminus, Inc.
 
 This program is free software: you can use, redistribute, and/or modify
 it under the terms of the GNU Affero General Public License, version 3
 or later ("AGPL"), as published by the Free Software Foundation.
 
 This program is distributed in the hope that it will be useful, but WITHOUT
 ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 FITNESS FOR A PARTICULAR PURPOSE.
 
 You should have received a copy of the GNU Affero General Public License
 along with this program. If not, see <http://www.gnu.org/licenses/>.
-->

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

