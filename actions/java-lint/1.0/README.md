### Java Action

用于 java 静态代码质量检查，采用阿里巴巴开源的代码检查插件 pmd 进行分析，制作成为 docker image 用于运行服务。

#### 使用

Examples:
1. pipeline.yml 增加 lint 描述
```yml
- java-lint:
    params:
      path: xxxx
```

2.在项目根 pom.xml 中增加如下 pmd plugin
```
<build>
    <plugins>
        <plugin>
            <groupId>org.apache.maven.plugins</groupId>
            <artifactId>maven-pmd-plugin</artifactId>
            <version>3.8</version>
            <configuration>
                <rulesets>
                    <ruleset>rulesets/java/ali-comment.xml</ruleset>
                    <ruleset>rulesets/java/ali-concurrent.xml</ruleset>
                    <ruleset>rulesets/java/ali-constant.xml</ruleset>
                    <ruleset>rulesets/java/ali-exception.xml</ruleset>
                    <ruleset>rulesets/java/ali-flowcontrol.xml</ruleset>
                    <ruleset>rulesets/java/ali-naming.xml</ruleset>
                    <ruleset>rulesets/java/ali-oop.xml</ruleset>
                    <ruleset>rulesets/java/ali-orm.xml</ruleset>
                    <ruleset>rulesets/java/ali-other.xml</ruleset>
                    <ruleset>rulesets/java/ali-set.xml</ruleset>
                </rulesets>
                <printFailingErrors>true</printFailingErrors>
            </configuration>
            <executions>
                <execution>
                    <goals>
                        <goal>check</goal>
                    </goals>
                </execution>
            </executions>
            <dependencies>
                <dependency>
                    <groupId>com.alibaba.p3c</groupId>
                    <artifactId>p3c-pmd</artifactId>
                    <version>1.3.0</version>
                </dependency>
            </dependencies>
        </plugin>
    </plugins>
</build>
```
