## mysql-assert
```yaml
version: "1.1"
stages:
  - stage:
      - mysql-assert:
          alias: mysql-assert
          description: 在对应的数据源中执行 sql 语句且可以断言和出参
          version: "1.0"
          params:
            datasource: xxxxx
            database: xxxx
            sql: show tables;
            out_params:
              - key: name
                expression: .[0].Tables_in_mysql
                value: columns_priv
                assert: =
              - key: name1
                expression: .[1].Tables_in_mysql
                value: db
                assert: =
```

### database
数据库名称

### datasource
数据源地址，图形界面可以下啦选择

### sql
执行的 sql 语句

### out_params

#### key 和 expression

out_params 下的 key 都会产生对应的出参，出参的值就是 expression 表达式从sql的返回值中解析出的值

假设 sql 查询结果如下： {"data":[{"id":185,"ss":"aa1"},{"id":186,"ss":"aa2"},{"id":187,"ss":"aa3"},{"id":188,"jsonStringValue":"{\"key\":\"value\"}"}]} </br>

使用 jq 表达式获取data下list中所有ID</br>
expression：[.data[] | try .id] </br>
结果为: [185,186,187,188]


使用 jq 表达式获取data下ss="aa2"的ID</br>
expression：.data[] | select(.ss=="aa2").id </br>
结果为：186

使用 jq 表达式获取data下ss="aa1"的ID，和data下ss="aa2"的ID，拼接成list[]</br>
expression：[(.data[] | select(.ss=="aa1").id),(.data[] | select(.ss=="aa2").id)]</br>
结果为：[185,186]

使用 jq 表达式获取data下标为3的字符串json转义后 json 中的值</br>
expression：.data[3].jsonStringValue | fromjson | .key </br>
结果为："value"


##### assert 和 value

assert 可图形下啦选择

上述解析出的值可以使用 assert 表达式和 value 进行断言，只要有一个断言失败，整个任务将会失败


