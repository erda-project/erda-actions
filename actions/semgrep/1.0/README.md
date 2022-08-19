### Semgrep Action 1.0

使用 semgrep 扫描代码, 快速进行代码分析，查找错误，在 CI 中运行安全扫描，并在整个组织中实施安全标准。

#### 最佳实践

##### 使用线上规则自动扫描
```yml
- stage:
  - semgrep:
      params:
        # 代码扫描目录
        code_dir: ${git-checkout}
        # 代码扫描规则，可以是本地文件或者线上规则
        config: auto
        # 报告输出格式
        format: sarif
        # 可选参数
        args:
          # 遇到错误时停止扫描
          - --error
```

##### 使用自定义规则扫描
以下是扫描是否注入了log4j的扫描规则
```yml
rules:
  - id: log4j-message-lookup-injection
    metadata:
      cwe: "CWE-74: Improper Neutralization of Special Elements in Output Used by a Downstream Component ('Injection')"
      owasp: 
        - "A03:2021 - Injection"
        - "A01:2017 - Injection"
      source-rule-url: https://www.lunasec.io/docs/blog/log4j-zero-day/
      references:
        - https://issues.apache.org/jira/browse/LOG4J2-3198
        - https://www.lunasec.io/docs/blog/log4j-zero-day/
        - https://logging.apache.org/log4j/2.x/manual/lookups.html
      category: security
      technology:
        - java
      confidence: MEDIUM
    message:
      Possible Lookup injection into Log4j messages. Lookups provide a way to add values to the Log4j messages at arbitrary
      places. If the message parameter contains an attacker controlled string, the attacker could inject arbitrary lookups,
      for instance '${java:runtime}'. This could lead to information disclosure or even remote code execution if 'log4j2.formatMsgNoLookups'
      is enabled. This was enabled by default until version 2.15.0.
    mode: taint
    pattern-sources:
    - patterns:
        - pattern: public $T $M(...)
    pattern-sinks:
    - patterns:
        - pattern: |
            (org.apache.log4j.Logger $L).$M(...)
    severity: WARNING
    languages:
      - java
```

假设将上面的规则放到代码根目录下并命名为`rules.yml`，则可以使用如下方式进行扫描
```yml
- stage:
  - semgrep:
      params:
        # 代码扫描目录
        code_dir: ${git-checkout}
        # 使用上问保存的log4j漏洞规则文件进行扫描
        config: rules.yml
        # 报告输出格式
        format: sarif
        # 可选参数
        args:
          # 遇到错误时停止扫描
          - --error
```

#### 可选参数
以下是常用的一些可选参数，更多参数可以参考 [官方文档](https://semgrep.dev/docs/cli-reference/#semgrep-scan-options)

1. `--error`: 遇到错误时停止扫描
2. `--use-git-ignore`: 跳过被git ignore忽略的文件
3. `--timeout`: 设置扫描超时时间，单位为秒
4. `--validate`: 校验配置的规则文件是否正确