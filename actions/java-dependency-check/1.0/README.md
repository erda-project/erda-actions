# Java-Dependency-Check Action

Java 依赖漏洞检测，并提供检测报告下载。

基本概念：

## CVE (Common Vulnerabilities and Exposures)
CVE 是已公开的信息安全漏洞字典，同时也是统一的漏洞编号标准。

CVE 可以简单的理解成对漏洞的收录和编号。

## CPE (Common Platform Enumeration)

CPE 是 通用平台枚举项，它是对 IT 产品的统一命名规范的标准化方法，包括系统、平台和软件包等，用于描述和识别企业计算资产中存在的应用程序、操作系统的硬件设备。

字段说明：
- part 目标类型，允许的值有 `a`（应用程序）、`h`（硬件平台）、`o`（操作系统）
- vendor 供应商
- product 产品名称
- version 版本号
- update 更新包
- edition 版本
- language 语言项

`内容格式：CPE:2.3:类型:厂商:产品:版本:更新版本:发行版本:界面语言:软件发行版本:目标软件:目标硬件:其他`

## CVSS (Common Vulnerability Scoring System)

CVSS 是工业标准的通用的漏洞评分系统，是用于描述安全漏洞严重程度的统一评分方案。

CVSS 是安全内容自动化协议（SCAP） 的一部分。通常 CVSS 与 CVE 一同由美国国家漏洞库（NVD）发布并保持数据的更新。

CVSS 的评分范围是 0-10.其中，10 是最高等级，不同机构按 CVSS 分支定义威胁的 高、中、低威胁级别。

#### 使用

Examples:

```yaml
- java-dependency-check:
    params:
      code_dir: ${git-checkout} # 代码目录
```

