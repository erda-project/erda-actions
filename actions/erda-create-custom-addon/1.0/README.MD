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

### Create Custom Addon
根据用户给定的名称和 Key-Value 创建 Custom Addon。如果同名 Addon 此前已被创建，则不会重复创建。

> 示例
```yaml
version: "1.1"
stages:
  - stage:
      - git-checkout:
  - stage:
      - custom-addon:
          params:
            name: java-demo
            tag: 默认配置
            conƒigs:
              key1: "01"
            template: ${{ dirs.git-checkout }}/custom-addon-tempate.yml
```

> custom-addon-template.yml
```yaml
name: java-demo
tag: 默认配置
default:
  key1: "31"
  key2: "32"
  key3: "33"
dev: 
  key1: "21"
  key2: "22"
test:
  key2: value2
staing: 
  key3: value3
prod:
  key4: value4
```

Key-Values 的优先级:
Action 配置 > 模板文件相应环境下的配置 > 模板文件下默认配置。
如以上示例中，创建开发环境下的 Custom Addon 时，`key1` 的值为 "01"，`key2` 的值为 "22"，
`key3` 的值为 "33"。

注意，value 的值都是字符串类型。
