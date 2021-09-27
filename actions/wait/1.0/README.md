### Wait Action

提供执行流水线时等待一段时间的功能

## 详细介绍

wait action 主要是支持用户在执行流水线时需要等待一段时间再执行下一个action的功能

## params

### wait_time_sec

必填。

等待的时间，单位为秒

## 使用

```yml
- wait:
    params:
      wait_time_sec: 300
```
