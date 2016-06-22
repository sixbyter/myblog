## 25 Nov 15 Mac加速docker下载hub.docker镜像

-  作者: [sixbyte](http://sixbyte.me/)
-  邮箱: liu.sixbyte@gmail.com
-  禁止转载

docker 提供了一个镜像仓库, [https://hub.docker.com](https://hub.docker.com), 里面有非常多优秀的镜像学习和使用, 但是 hub.docker.com 在国内的访问速度非常非常糟糕, 所以我希望使用代理的方式进行访问

docker 官方允许在创建 docker machine 的时候设置 `http_proxy`, 但是我没有静态的 `http_proxy` 服务器, 所以无法使用.

在网上找到了 `polipo` 和 `shadowsocks` 的结合可以解决这个问题, 正好我也是在使用 `shadowsocks`, 所以写了这个办法

### 要求

brew, shadowsocks

### 通过brew安装polipo

如果出现权限问题, 执行`sudo mkdir /usr/local/var && sudo chown -R "$USER":admin /usr/local/var`

```shell
brew install polipo
```

### socks 转 http_proxy

`192.168.1.147` 为我的本机ip

```shell
polipo socksParentProxy=localhost:1080 proxyAddress="192.168.1.147"
```

### 执行命令通过http_proxy

```
http_proxy=http://192.168.1.147:8123 composer self-update
http_proxy=http://192.168.1.147:8123 curl google.com
```

### docker-machine 内部 命令通过 http_proxy 执行命令

`192.168.1.147` 为我的本机ip

```shell
http_proxy=http://192.168.1.147:8123 docker pull ubuntu
```
