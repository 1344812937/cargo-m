# cargo-m
- 将本地Maven仓库分享问网络仓库, 自动扫描本地仓库资源。 首次运行将自动生成配置文件，可修改配置文件后重启。
```toml
[http]
  host = '监听IP'
  port = '网络服务端口'

[maven_repo]
  enabled = true  ＃是否启用自动扫描
  local_path = '本地资源路径'
```
- 服务启动成功后可将远程仓库地址设置为一下内容, 可直接拉取依赖。
```azure
http://ip:port/maven-repo/getRepo/
```