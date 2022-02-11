### erda-get-addon-info Action

查看指定 runtime 的 指定名称的 addon 的配置信息。

#### 使用

```yml
- erda-get-addon-info:
    params:
      application_name: ${application_name}
      runtime_id: ${runtime_id}
      addon_name: ${addon_name}
```
**注意**:
* 参数 addon_name 必须提供
* 参数 application_name 与 runtime_id 至少提供一个。

执行结果是输出一组 Meta 表示获取的 addon 的配置信息，对应的 Meta 的名称是配置项的 Key，value 是配置项的值

后续 Action 可以通过 ${{ outputs.alias.val }} 获取并使用 Meta 信息。

例如要获取 addon 配置 MYSQL_HOST 的地址，则后续 Action 可以通过 ${{ outputs.MYSQL_HOST.val }} 获取并使用 Meta 信息。