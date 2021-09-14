## jsonparse


假设 data 结果如下： {"data":[{"id":185,"ss":"aa1"},{"id":186,"ss":"aa2"},{"id":187,"ss":"aa3"},{"id":188,"jsonStringValue":"{\"key\":\"value\"}"}]} </br>

   使用 jq 表达式获取data下list中所有ID</br>
   expression：[.data[] | try .id] </br>
   结果为: [185,186,187,188]
   
```yaml
- jsonparse:
    alias: jq
    version: "1.0"
    params:
      out_params:
        - key: "name"
          expression: "[.data[] | try .id]"
      data: '{"data":[{"id":185,"ss":"aa1"},{"id":186,"ss":"aa2"},{"id":187,"ss":"aa3"},{"id":188,"jsonStringValue":"{\"key\":\"value\"}"}]}'
```   
   
   使用 jq 表达式获取data下ss="aa2"的ID</br>
   expression：.data[] | select(.ss=="aa2").id </br>
   结果为：186

```yaml
- jsonparse:
    alias: jq
    version: "1.0"
    params:
      out_params:
        - key: "name"
          expression: '.data[] | select(.ss=="aa2").id'
      data: '{"data":[{"id":185,"ss":"aa1"},{"id":186,"ss":"aa2"},{"id":187,"ss":"aa3"},{"id":188,"jsonStringValue":"{\"key\":\"value\"}"}]}'
```      
   
   使用 jq 表达式获取data下ss="aa1"的ID，和data下ss="aa2"的ID，拼接成list[]</br>
   expression：[(.data[] | select(.ss=="aa1").id),(.data[] | select(.ss=="aa2").id)]</br>
   结果为：[185,186]

```yaml
- jsonparse:
    alias: jq
    version: "1.0"
    params:
      out_params:
        - key: "name"
          expression: '[(.data[] | select(.ss=="aa1").id),(.data[] | select(.ss=="aa2").id)]'
      data: '{"data":[{"id":185,"ss":"aa1"},{"id":186,"ss":"aa2"},{"id":187,"ss":"aa3"},{"id":188,"jsonStringValue":"{\"key\":\"value\"}"}]}'
```      
   
   使用 jq 表达式获取data下标为3的字符串json转义后 json 中的值</br>
   expression：.data[3].jsonStringValue | fromjson | .key </br>
   结果为："value"   

```yaml
- jsonparse:
    alias: jq
    version: "1.0"
    params:
      out_params:
        - key: "name"
          expression: ".data[3].jsonStringValue | fromjson | .key"
      data: '{"data":[{"id":185,"ss":"aa1"},{"id":186,"ss":"aa2"},{"id":187,"ss":"aa3"},{"id":188,"jsonStringValue":"{\"key\":\"value\"}"}]}'
```      


  取多个值

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
  使用 jackson 取值

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

  使用jsonparse解析原生取值

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

  使用其他 action 的出参放入 data 中进行取值

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


