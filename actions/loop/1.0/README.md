### Loop Action

循环某个操作，直至满足退出条件或超时。

#### 使用

Examples:

1. HTTP

```yaml
- loop:
    params:
      type: HTTP
      json_response_success_field: .success
      loop_max_times: 10
      loop_interval: 1s
      loop_decline_ratio: 2
      loop_decline_limit: 10s
      http_conf:
        host: 127.0.0.1:8080
        path: /ping
        method: GET
        header:
          OPENAPI_TOKEN: xxx
          Org-ID: 1
        query_params:
          k1: v1
          k2: v2
        response_body: '{"key":"value"}'
```

2. CMD

```yaml
- loop:
    params:
      type: CMD
      json_response_success_field: .success
      loop_max_times: 10
      loop_interval: 1s
      loop_decline_ratio: 2
      loop_decline_limit: 10s
      CMD: curl 127.0.0.1:80/ping
```
