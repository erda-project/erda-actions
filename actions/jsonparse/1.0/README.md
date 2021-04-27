## jsonparse

根据 out_params 中

```yaml
- jsonparse:
    alias: jq
    version: "1.0"
    params:
      out_params:
        - key: "name"
          expression: ".name"
        - key: "id"
          expression: ".id"
      data: '{"name": 123, "id": 123}'
```

```yaml
- jsonparse:
    alias: jackson
    version: "1.0"
    params:
      out_params:
        - key: "name"
          expression: "$.name"
        - key: "id"
          expression: "$.id"
      data: '{"name": 123, "id": 123}'
```

```yaml
- jsonparse:
    alias: jsonparse
    version: "1.0"
    params:
      out_params:
        - key: "name"
          expression: "name"
        - key: "id"
          expression: "id"
      data: '{"name": 123, "id": 123}'
```

```yaml
- jsonparse:
    alias: otherActionOutput
    version: "1.0"
    params:
      out_params:
        - key: "name"
          expression: ".name"
        - key: "id"
          expression: ".id"
      data: ${{ outputs.actionName.actonOuput }}  # other action output
```


