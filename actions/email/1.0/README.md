## email 发送

### 模板文件(例子)

```text
{{range .email_template_object}}
    <p>{{ .test }}</p>
{{end}}
```

### 使用代码仓库存储模板

代码仓库 -> 代码预览 -> 根目录下创建一个新文件，将上面的内容粘贴进去，文件命假定为 `template_email.txt`

```yaml
version: "1.1"
stages:
  - stage:
      - git-checkout:
          alias: git-checkout
          description: 代码仓库克隆
          version: "1.0"
  - stage:
      - email:
          alias: email
          description: 邮件发送
          version: "1.0"
          params:
            email_template_addr: ${{ dirs.git-checkout }}/template_email.txt
            email_template_object: '[{"test":1},{"test":2}]'
            to_email:
              - xxx@xxx.com
```

### 使用流水线变量存储模板

应用设置 -> 流水线 -> 变量配置 -> 选择文件类型，将上面的文件内容上传，key 设定为例如 `template_email`

```yaml
version: "1.1"
stages:
  - stage:
      - email:
          alias: email
          description: 邮件发送
          version: "1.0"
          params:
            email_template_addr: ${{ configs.template_email }}
            email_template_object: '[{"test":1},{"test":2}]'
            to_email:
              - xxx@xxx.com
```

### 说明

#### 参数说明
email_template_addr: 模板地址，目前只支持 erda 的变量设置和流水线中的文件地址，网络文件暂时不支持
email_template_object: 模板文件中占位符渲染所使用的 json 结构体
to_email: 接收邮件的邮箱地址，数组结构，可以同时发送多个人

#### 模板说明
模板渲染是使用 golang 语言的官方的 template 库，具体语法参考 `https://pkg.go.dev/text/template`

#### 发送方 smtp 服务器说明

默认邮件发送使用的是中心集群配置的 smtp 地址，用户如果想要使用自己的 smtp 服务器可以在 应用设置 -> 流水线 -> 变量配置 中配置下面几个值，或者修改中心集群的 smtp 配置来全局修改

SMTP_HOST: 邮件服务器的地址
SMTP_PORT: 邮件服务器的端口
SMTP_EMAIL: 邮件服务器的发送地址
SMTP_PASSWORD: 邮件服务器的密码
