# Contrast Security Action

## 简介
ContrastSecurity用于安全漏洞扫描。

## 额外参数说明

### assert_count
用户断言安全漏洞数量，默认为0，代表不进行断言，当填写的assert_count大于0时，会判断API返回的漏洞总数是否大于assert_count。


#### 示例
```yml
  - contrast-security:
      alias: contrast-security
      description: contrast安全扫描
      version: "1.0"
      params:
        api_key: xxx
        app_id: xxx
        org_id: xxx
        service_key: xxx
        username: xxx
```