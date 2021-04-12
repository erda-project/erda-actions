# OSS Action

上传本地文件或目录到OSS

## 使用

```yml
- oss-upload:
    params:
      access_key_id: ((oss.key))
      access_key_secret: ((oss.secret))
      endpoint: http://oss-cn-hangzhou.aliyuncs.com
      bucket: terminus-dice
      local_path: dist
      oss_path: web
```
