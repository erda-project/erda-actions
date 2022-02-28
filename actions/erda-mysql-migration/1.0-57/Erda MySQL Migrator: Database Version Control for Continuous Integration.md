<a name="ffENi"></a>
# 为什么要进行数据库版本控制？
现代软件工程逐渐向持续集成、持续交付演进已成时代趋势，我们的软件不是一次性交付了事，而是要持续地集成和交付，
如何高效地进行软件版本控制成为我们面临的挑战；我们的软件也不是只要部署到某一套环境中，
而是需要部署到开发、测试、生产以及更多的客户环境中，如何一套代码适应不同的环境成为我们要思考的问题。

![img.png](http://terminus-paas.oss-cn-hangzhou.aliyuncs.com/paas-doc/2022/02/28/15fc7a58-e8ab-4d4f-b257-d96667dc5ed9.png)

一套软件的副本要部署在不同的环境（图源：Flyway）<br />代码版本管理工具（Git、SVN 等）和托管平台（Github、Erda DevOps Platform 等）让我们能有效地进行代码版本管理。越来越丰富的 CI/CD 工具让我们能定义可重复的构建和持续集成流程，发布和部署变得简单清晰。<br />“基础设施即代码”的思想，让我们可以用代码定义基础设施，从而抹平了各个环境的差异。 <br />可以说，在软件侧我们应对这些挑战已经得心应手。<br />
但是绝大多数项目都至少包含两个重要部分：业务软件，以及业务软件所使用的数据库——许多项目数据库侧的版本控制仍面临乱局：

- 很多项目的数据库版本控制仍依赖于“人肉维护”，需要开发者手动执行 SQL
- 环境一多，几乎没人搞得清某个环境上数据库是什么状态了
- database migrations 脚本没有统一管理起来，遗失错漏时有发生
- 不知道脚本的状态是应用了呢还是没有应用呢，也许在这个环境应用了但在那个环境却没有应用 ？
- 脚本里有一行破坏性代码，执行了后将一个表字段删除了，数据恢复不来，只能“从删库到跑路”
- ……

为了应对这样的乱局，我们需要数据库版本控制工具。数据库版本控制，即 Database Migration，它能帮你

- 管理数据库的定义和迁移历程
- 在任意时刻和环境从头创建数据库至指定的版本
- 以确定性的、安全的方式执行迁移
- 清楚任意环境数据库处于什么状态

从而让数据库与软件的版本管理同步起来，软件版本始终能对应到正确的数据库版本，同时提高安全性、降低维护成本。
<a name="aKH38"></a>
# Erda 项目是如何做 database migration 的？
Erda 是基于多云架构的一站式企业数字化平台，为企业提供 DevOps、微服务治理、多云管理以及快数据管理等云厂商无绑定的 IT 服务。Erda 既可以私有化交付，也提供了 SaaS 化云平台 [Erda Cloud](https://www.erda.cloud)，以及开源的社区版。<br />当你正在阅读这篇文章时，有无数来自不同组织的应用程序正在 Erda Cloud 或 Erda 私有化平台的流水线上完成以构建和部署为核心的 CI/CD 流程，无数的代码，以这种持续而自动化的方式转化成服务实例。<br />Erda 平台不但接管了这些组织的应用程序的集成、交付，Erda 项目自身的集成也是托管在 Erda DevOps 平台的。Erda 自身的持续集成和丰富的交付场景要求它能进行安全、高效、可持续的数据库版本控制，托管在 Erda 上的应用程序也要求 Erda 提供一套完整的数据库版本控制方案。<br />Erda 项目使用 _Erda MySQL Migrator_ 作为数据库版本控制工具，它被广泛应用于 CI/CD 流程和命令行工具中。
<a name="PXyHb"></a>
## 基本原理
第一次使用 Erda MySQL Migrator 进行数据库版本控制时会在数据库中新建一个名为 _schema_migration_history _的表。

![](http://terminus-paas.oss-cn-hangzhou.aliyuncs.com/paas-doc/2022/02/28/5579adf7-abd3-42a9-9c26-7a726bc06b8b.png)

schema_migration_history 表的基本结构(部分主要字段)<br />Erda MySQL Migrator 每次执行 database migration 时，会对比文件系统中的 migrations 脚本和 schem_migration_history 表中的执行记录，标记出增量的部分。在一系列审查后，Erda MySQL Migrator 将增量的部分应用到目标 database 中。成功应用的脚本被记录在案。
<a name="UAWqT"></a>
## 使用 Erda MySQL Migrator 命令行工具进行数据库版本控制
<a name="NmXXX"></a>
### erda-cli 工具的安装与使用
`erda-cli` 是 erda 项目命令行工具，它集成了 Erda 平台安装、Erda 拓展管理以及开发脚手架。其中 `erda-cli migrate` 命令集成了数据库版本控制全部功能。<br />从 [erda 仓库](https://github.com/erda-project/erda) 拉取代码到本地，切换到 master 分支，执行以下命令可以编译`erda-cli`
```shell
% make prepare-cli
% make cli
```
注意编译前应确保当前环境已安装 docker。编译成功后项目目录下生成了一个 `bin/erda-cli` 可执行文件。
<a name="NctRX"></a>
### 使用 erda-cli migrate 进行数据库版本迁移
Erda MySQL Migrator 要求按 `modules/scripts` 两级目录组织数据库版本迁移脚本，以 erda 仓库为例：
```shell
.erda/migrations
├── apim
│   ├── 20210528-apim-base.sql
│   ├── 20210709-01-add-api-example.py
│   └── requirements.txt
... ...
├── cmdb
│   ├── 20210528-cmdb-base.sql
│   ├── 20210610-addIssueSubscriber.sql
│   ├── 20210702-updateMbox.sql
│   └── 20210708-add-manageconfig.sql
└── config.yml
    └── 20200528-tmc-base.sql

```
erda 项目将数据库迁移脚本放在 `.erda/migrations` 目录下，目录下一层级是按模块名（微服务名）命名的脚本目录，其各自下辖本模块所有脚本。与脚本目录同级的，还有一个 config.yml 的文件，它是 Erda MySQL Migration 规约配置文件，它描述了 migrations 脚本所需遵循的规约。<br />脚本目录下按文件名字符序排列着 migrations 脚本，目前支持 SQL 脚本和 Python 脚本。如果目录下存在 Python 脚本，则需要用 `requirements.txt` 来描述 Python 脚本的依赖。<br />我们先一起来看下 `erda-cli migrate` 命令的帮助信息：
```shell
% erda-cli migrate --help

erda-cli migrate --mysql-host localhost --mysql-username root --mysql-password my_password --database erda

Usage:
  erda-cli migrate  [flags]
  erda-cli migrate [command]

Examples:
erda-cli migrate --mysql-host localhost --mysql-username root --mysql-password my_password --database erda

Available Commands:
  lint        Erda MySQL Migration lint
  mkpy        make a python migration script pattern
  record      manually insert the migration record

Flags:
      --database string         [MySQL] database to use. env: ERDA_MYSQL_DATABASE
      --debug-sql               [Migrate] print SQLs
  -h, --help                    help for migrate
      --lint-config string      [Lint] Erda MySQL Lint config file
      --modules strings         [Lint] the modules for migrating
      --mysql-host string       [MySQL] connect to host. env: ERDA_MYSQL_HOST
      --mysql-password string   [MySQL] password to use then connecting to server. env: ERDA_MYSQL_PASSWORD
      --mysql-port int          [MySQL] port number to use for connection. env: ERDA_MYSQL_PORT (default 3306)
      --mysql-username string   [MySQl] user for login. env: ERDA_MYSQL_USERNAME
      --output string           [Migrate] the directory for collecting SQLs
      --sandbox-port int        [Sandbox] sandbox expose port. env: ERDA_SANDBOX_PORT (default 3306)
      --skip-lint               [Lint] don't do Erda MySQL Lint
      --skip-mig                [Migrate] skip doing pre-migration and real migration
      --skip-sandbox            [Migrate] skip doing migration in sandbox

```
erda-cli migrate 命令行的帮助信息<br />参数解释：<br />--mysql-host, --mysql-port, --mysql-username, --mysql-password, --database 等参数是连接目标 MySQL Server 所需的参数。<br />--debug-sql 决定是否打印执行的 SQL 语句到标准输出。<br />--lint-config 是指 Erda MySQL Migration 规约配置文件的路径，如不设置，则使用默认配置。<br />--modules 决定执行哪些模块下的数据库版本迁移，默认情况下执行所有模块的迁移。<br />--output 是 SQL 执行日志输出目录。<br />--sandbox-port 是 MySQL 沙盒暴露的端口。<br />--skip-lint 表示跳过 Erda MySQL Migration 规约检查，默认为 false。<br />--skip-mig 跳过在 MySQL Server 上执行 migration，相当于 Dryrun。<br />--skip-sandbox 跳过在 MySQL 沙盒中执行 migration。<br />​

进入 migrations 脚本所在目录 `.erda/migrations`，执行 `erda-cli migrate`：
```shell
% erda-cli migrate --mysql-host localhost \
    --mysql-username root \
    --mysql-password 123456789 \
    --sandbox-port 3307 \
    --database erda
INFO[0000] Erda Migrator is working                     
INFO[0000] DO ERDA MYSQL LINT....                       
INFO[0000] OK                           
INFO[0000] DO FILE NAMING LINT....                        
INFO[0000] OK                            
INFO[0000] DO ALTER PERMISSION LINT....                 
INFO[0000] OK                     
INFO[0000] DO INSTALLED CHANGES LINT....                
INFO[0000] OK                    
INFO[0000] COPY CURRENT DATABASE STRUCTURE TO SANDBOX.... 
INFO[0014] OK 
INFO[0014] DO MIGRATION IN SANDBOX....                  
INFO[0014] OK                                            
INFO[0014] DO MIGRATION....                             
INFO[0014]                 module=apim
... ...
INFO[0014]                 module=cmdb
INFO[0014] OK
INFO[0014] Erda MySQL Migrate Success !
```
执行 erda-cli migrate 命令<br />从执行日志可以看到，命令行执行一系列检查以及沙盒预演后，成功应用了本次 database migration。我们可以登录数据库查看到脚本的应用情况。
```sql
mysql> SELECT service_name, filename FROM schema_migration_history;
+---------------+-------------------------------------------+
| service_name  | filename                                  |
+---------------+-------------------------------------------+
| apim          | 20210528-apim-base.sql                    |
| apim          | 20210709-01-add-api-example.py            |
... ...				... ...																			... ...
| cmdb          | 20210528-cmdb-base.sql                    |
| cmdb          | 20210610-addIssueSubscriber.sql           |
| cmdb          | 20210702-updateMbox.sql                   |
| cmdb          | 20210708-add-manageconfig.sql             |
+---------------+-------------------------------------------+
```
登录 MySQl Server 查看脚本应用情况
<a name="yToPF"></a>
### 基于 Python 脚本的 data migration
从上一节我们看到，脚本目录中混合着 SQL 脚本和 Python 脚本，migrator 对它们一致地执行。**Erda MySQL Migrator 在设计之初就决定了单脚本化的 migration，即一个脚本表示一次 migration 过程。**大部分 database migration 都可以很好地用 SQL 脚本表达，但仍有些包含复杂逻辑的 data migration 用 SQL 表达则会比较困难。对这类包含复杂业务逻辑的 data migration，Erda MySQL Migrator 支持开发者使用 Python 脚本。<br />erda-cli 提供了一个命令行`erda-cli migrate mkpy`来帮助开发者创建一个基础的 Python 脚本。
```shell
% erda-cli migrate mkpy --help
make a python migration scritp pattern.

Usage:
  erda-cli migrate mkpy  [flags]

Examples:
erda-cli migrate mkpy --module my_module --name my_script_name

Flags:
  -h, --help             help for mkpy
  -m, --module string    migration module name
  -n, --name string      script name
      --tables strings   dependency tables
      --workdir string   workdir (default ".")
```
erda-cli migrate mkpy 的 help 信息<br />参数说明：<br />--module：脚本所在的模块名。<br />--name：脚本名称。<br />--tables：脚本中要引用到的数据库表名。<br />--workdir：migrations 目录。<br />​

执行
```shell
% erda-cli migrate mkpy --module my_module --name myfeature.py --tables blog,author,info
```
命令生成如下脚本：
```python
"""
Generated by Erda Migrator.
Please implement the function entry, and add it to the list entries.
"""


import django.db.models


class Blog(django.db.models.Model):
    name = models.CharField(max_length=100)
    tagline = models.TextField()

    class Meta:
        db_table = "blog"

class Author(django.db.models.Model):
    name = models.CharField(max_length=200)
    email = models.EmailField()

    class Meta:
        db_table = "author"

class Info(django.db.models.Model):
    blog = models.ForeignKey(Blog, on_delete=models.CASCADE)
    headline = models.CharField(max_length=255)
    body_text = models.TextField()
    pub_date = models.DateField()
    mod_date = models.DateField()
    authors = models.ManyToManyField(Author)
    number_of_comments = models.IntegerField()
    number_of_pingbacks = models.IntegerField()
    rating = models.IntegerField()

    class Meta:
        db_table = "info"


def entry():
    """
    please implement this and add it to the list entries
    """
    pass


entries: [callable] = [
    entry,
]
```
该脚本可以分为四个部分：<br />1、import 导入必要的包。脚本中采用继承了 django.db.models.Model 的类来定义库表，因此需要导入 django.db.model 库。开发者可以根据实际情况导入自己所须的包，但由于单脚本提交的原则，脚本中不应当导入本地其他文件。<br />2、模型定义。脚本中 `class Blog`、 `class Author` 和 `class Entry`是命令行工具为开发者生成的模型类。开发者可以使用命令行参数 `--tables` 指定要生成哪些模型定义，以便在开发中引用它们。注意，生成这些模型定义类时并没有连接数据库，而是根据文件系统下过往的 migration 所表达的 Schema 生成。生成的模型定义只表示了表结构而不包含表关系，如“一对一”、“一对多”、“多对多”等。如果开发者要使用关联查询，应当编辑模型，自行完成模型关系的描述。Django ORM 的模型关系仅表示逻辑层面的关系，与数据库物理层的关系无关。<br />3、entry 函数。命令行为开发者生成了一个名为 `entry` 的函数，但是没有任何函数体。开发者须要自行实现该函数体以进行 data migration。<br />4、entries，一个以函数为元素的列表，是程序执行的入口。开发者要将实现 data migration 的业务函数放到这里，只有 entries 中列举的函数才会被执行。<br />从以上脚本结构可以看到，我们选用的 Django ORM 来描述模型和进行 CRUD 操作。为什么采用 Django ORM 呢？因为 Django 是 Python 语言里最流行的 web 框架之一，Django ORM 也是 Python 中最流行的 ORM 之一，其设计完善、易用、便于二次开发，且有详尽的文档、丰富的学习材料以及活跃的社区。无论是 Go 开发者还是 Java 开发者，都能在掌握一定的 Python 基础后快速上手该 ORM。我们通过两个简单的例子来了解下如何利用 Django ORM 来进行 CRUD 操作。<br />示例 1：创建一条新记录。
```python
# 示例 1
# 创建一条记录
def create_a_blog():
    blog = Blog()
    blog.name = "this is my first blog"
    blog.tagline = "this is looong text"
    blog.save()
```
Django ORM 创建一条记录十分简单，引用模型类的实例，填写字段的值，调用 `save()`方法即可。<br />示例 2：删除所有标题中包含 "Lennon" 的 Blog 条目。Django 提供了一种强大而直观的方式来“追踪”查询中的关系，在幕后自动处理 SQL **JOIN** 关系。它允许你跨模型使用关联字段名，字段名由双下划线分割，直到拿到想要的字段。
```python
# 示例 2
# 删除所有标题中含有 'Lennon' 的 Blog 条目:
def delete_blogs_with_headline_lennon():
    Blog.objects.filter(info__headline__contains='Lennon').delete()
```
最后，别忘了将这两个函数放到 entries 列表中，不然它们不会被执行。
```python
entries: [callable] = [
    create_a_blog,
    delete_blogs_with_headline_lennon,
]
```
可以看到，编写基于 Python 的 data migration 是十分方便的。`erda-cli migrate mkpy` 命令行为开发者生成了模型定义，引用模型类及其实例可以便捷地操作数据变更，开发只须关心编写函数中的业务逻辑。<br />​

进一步了解 Django ORM 的使用请查看文档：<br />

<a name="NwMLU"></a>
## 使用 Erda MySQL 数据迁移 Action 在 CI/CD 时进行数据库版本控制
每日凌晨，Erda 上的一条流水线静静启动，erda 仓库的主干分支代码都会被集成、构建、部署到集成测试环境。开发者一早打开电脑，
登录集成测试环境的 Erda 平台验证昨日集成的新 feature 是否正确，发现昨天新合并的 migrations 也一并应用到了集成测试环境。
这是怎么做到的呢 ？

![](http://terminus-paas.oss-cn-hangzhou.aliyuncs.com/paas-doc/2022/02/28/e23243aa-8ce2-46b0-8161-ca727983c74d.png)

Erda 每日自动化集成流水线（部分步骤）<br />原来这条流水线每日凌晨拉取 erda 仓库主干分支代码 -> 构建应用 -> 将构建产物制成部署制品 -> 在集成测试环境执行 Erda MySQL 数据迁移 -> 将制品部署到集成测试环境。**流水线中的 _Erda MySQL 数据迁移 _节点是集成了 Erda MySQL Migrator 全部功能的 Action，是 Erda MySQL Migrator 在 Erda CI/CD 流水线中的应用。**<br />Erda MySQL Migrator 除了可以作为 Action 编排在流水线中，还可以脱离 Erda 平台作为命令行工具单独使用。
<a name="oEUWa"></a>
# Erda MySQL Migrator 其他特性
<a name="AeVHS"></a>
## 规约检查
Erda 团队为了团队协作的统一规范制定了一系列开发规约，其中包含《Erda MySQL Migration 规约》。Erda MySQL Migrator 工具可以帮助开发者检查提交的 migration 脚本是否符合 《Erda MySQL Migration 规约》。
<a name="AMBRf"></a>
### 使用命令行工具进行规约检查
`erda-cli migrate lint` 命令可以检查指定目录下所有脚本的 SQL 语句是否符合规约。 开发者在编写 migration 时用该命令来预先检查，避免提交不合规了不合规的脚本。<br />例如开发者在 SQL 脚本中编写了如下语句：
```sql
alter table dice_api_assets add column col_name varchar(255);
```
执行规约检查：
```shell
% erda-cli migrate lint

2021/07/19 17:39:43 Erda MySQL Lint the input file or directory: .
apim/20210715-01-feature.sql:
    dice_api_assets:
        - 'missing necessary column definition option: COMMENT'
        - 'missing necessary column definition option: NOT NULL'

apim/20210715-01-feature.sql [lints]
apim/20210715-01-feature.sql:1: missing necessary column definition option: COMMENT: 
~~~> alter table dice_api_assets add column col_name varchar(255);

apim/20210715-01-feature.sql:1: missing necessary column definition option: NOT NULL: 
~~~> alter table dice_api_assets add column col_name varchar(255);

```
使用命令行工具进行本地规约检查<br />可以看到命令行返回了检查报告，指出了某个文件中存在不合规的语句，并指出了具体的文件、行号、错误原因等信息。上面示例中指出了这条语句有两条不合规处：一是新增列时，应当有列注释，此处缺失；二是新增的列应当是 NOT NULL 的，此处没有指定。

<a name="PNPn3"></a>
### 使用 Erda MySQL Migration Action 进行规约检查
略。

<a name="EGlKI"></a>
## 沙盒与 Dryrun
引入沙盒是为了在将 migrations 应用到目标数据库前进行一次模拟预演，期望将问题的发现提前，防止将问题 migration 应用到了目标数据库中。
<a name="anFmM"></a>
## 文件篡改检查与修订机制
Erda MySQL Migrator 不允许篡改已应用过的文件。之所以这样设计是因为一旦修改了已应用过的脚本，那么代码与真实数据库状态就不一致了。如果要修改表结构，应当增量地提交新的 migrations。这是一种常见的做法，Flyway 等工具也会对已执行的文件进行检查。但实际生产中，“绝不修改过往文件”这种理想状态很难达到，Erda MySQL Migrator 提供了一种修订机制。当用户想修改一个文件名为“some-feature.sql”过往文件时，他应该修改该文件，并提交一个名为“patch-some-feature.sql”的包含了修改内容的文件到 .patch 目录中。
<a name="DBQMz"></a>
## 日志收集
Erda MySQL Migrator 在 debug 模式下，会打印所有执行执行过程和 SQL 的标准输出。除此之外，它还可以将纯 SQL 输出到指定目录的日志文件中。
<a name="VscTn"></a>
## 基线准确性检查
Erda MySQL Migrator 将第一次小娟说数据库版本控制时线上数据库已有的库表状态视为“基线”状态。基线准确性是指，文件系统中的脚本描述的库表状态与线上数据库真实状态的一致性。Flyway 等工具也引入了“基线”的概念，这些工具大多没有对基线脚本准确性进行判断。然而这种判断是十分重要的，当基线脚本描述的库表结构与真实的数据库结构不一致而我们却把它们当作正确的基线记录下来，会产生严重的隐患。例如数据库中某个表字段类型是 VARCAHR(255)，而基线文件中是 VARCHAR(191)，这样细微的差别在大部分情况下都不会引发问题，然而当某一天我们在这个字段上建立索引时，就会因索引超长而造成问题。<br />因此 Erda MySQL Migrator 在中途介入数据库版本控制时会对用户整理的基线准确性进行检查，如果基线描述的库表结构与目标数据库已存在的库表结构不一致，则会拒绝继续执行。
<a name="Hc4LJ"></a>
# 获取工具
<a name="otqid"></a>
### erda-cli 下载地址
Mac：[http://erda-release.oss-cn-hangzhou.aliyuncs.com/cli/mac/erda-cli](http://erda-release.oss-cn-hangzhou.aliyuncs.com/cli/mac/erda-cli)<br />Linux：[http://erda-release.oss-cn-hangzhou.aliyuncs.com/cli/linux/erda-cli](http://erda-release.oss-cn-hangzhou.aliyuncs.com/cli/mac/erda-cli)<br />注意：以上 erda-cli 仅用于 amd64 平台，其他平台请按文中介绍的安装方式自行构建。
<a name="WyIsi"></a>
### 服务市场（Actions）
Erda MySQL MIgrator Action 源码地址：[https://github.com/erda-project/erda-actions/tree/master/actions/erda-mysql-migration/1.0-57](https://github.com/erda-project/erda-actions/tree/master/actions/erda-mysql-migration/1.0-57)

服务市场：[https://www.erda.cloud/market/action/erda-mysql-migration](https://www.erda.cloud/market/action/erda-mysql-migration)<br />
