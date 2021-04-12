### python-script Action

### 一、简介

执行用户的 python 的脚本

### 二、详细介绍
用户有一段自己写的 python 脚本，通过 git-checkout 或者其他方式注入到本action中
action拿到脚本路径，执行用户的python脚本

### 三、使用场景
用户有一段代码存在与 gittar 或者其他存储中，通过一个上游action拉取脚本，在本python-script脚本中引用上游路径
python2 执行python XX.py 执行用户的脚本
python3 执行python3 XX.py 执行用户的脚本
### 四、使用方式
```yaml
version: "1.1"
stages:
  - stage:
      - git-checkout:
          params:
            depth: 1
  - stage:
      - python-script:
          params:
            commands:
            - cd /
            - python3 ${git-checkout}/test.py
```