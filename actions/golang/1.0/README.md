## GoLang Action

### 额外参数说明

#### assets

assets用于配置文件等资源文件、

```
只写一个文件名则保持一样的目录结构
srs:dest 可以指定目标文件的路径
```

#### goproxy

用于 go mod模式下的代理配置


#### 示例

如果一个项目中只有一个go模块并且main函数在顶级包名中级不需要额外配置command和target参数，可以自动探测

```yml
  - golang:
      params:
        context: ${repo}
        service: hello
```


```yml
  - golang:
      params:
        command: go build -o hello cmd/main.go
        target: hello
        context: ${repo}
        service: hello
        assets:
            - application.yml
            - log.json:conf/log.json
```

