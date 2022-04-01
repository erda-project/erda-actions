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

### Project Artifacts
> ### 项目打包发布制品
将给定的应用的制品组合发布为项目制品

![](http://terminus-paas.oss-cn-hangzhou.aliyuncs.com/paas-doc/2022/03/02/a43327bb-5989-477d-ad84-5e31893cebcd.png)
> 示例
```yaml
version: "1.1"
stages:
  - stage:
      - project-artifacts:
          alias: project-artifacts
          description: 应用制品发布到项目制品
          version: "1.0"
          params:
            changeLog: auto compose from applications // 项目制品的内容
            groups:                                   // 应用制品分组，有一个应用描述列表组成
              - applications:                         // 应用列表
                  - branch: release/1.0               // 应用 release 分支, 选择该分支下最新的 release 发布到项目
                    name: go-demo                     // 应用名称
                  - name: java-demo
                    releaseID: a9af810ebd884107a3b9a  // 指定的应用 releaseID, 优先级高于 branch
            modes:                                    // 部署模式，优先级高于 groups
              default: 
                dependOn:                             // 依赖模式
                  - modeA
                expose: true                          // 是否展示
                groups:
                  - applications:                         
                      - branch: release/1.0               
                        name: go-demo                     
                      - name: java-demo
                        releaseID: a9af810ebd884107a3b9a
              modeA:
                expose: false
                groups:
                  - applications:
                      - branch: release/1.0
                        name: base

            version: 1.0.0+${{ random.date }}         // 项目 release 的版本号
```
