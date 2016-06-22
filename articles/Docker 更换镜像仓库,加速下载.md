## 22 Dec 15 Docker 更换镜像仓库,加速下载

-  作者: [sixbyte](http://sixbyte.me/)
-  邮箱: liu.sixbyte@gmail.com
-  禁止转载


之前写了一篇文章, [Mac加速docker下载hub.docker镜像](), 但是!! 效果相当不理想, 对于 `curl` 等其他纯 `http` 请求的命令效果是想多不错. 但是 `docker` pull 的速度已经触碰到了我的底线.
在用 `docker` 部署 `gitlab` 时, 我使用的是 [sameersbn/docker-gitlab](https://github.com/sameersbn/docker-gitlab) 的镜像, 整个安装过程其实非常简单, 只需要执行2条命令, 然而, 却在下载镜像花费了一个晚上还没完成, 链接 `docker hub` 的异常的慢, 经常出现 `timeout` :(

为此, 我不得不考虑国内的镜像加速服务. 目前在 `google` 上搜索到比较多的资料是 [DaoCloud](http://www.daocloud.io/)

### 更换镜像仓库

我一直不肯使用国内镜像加速服务的原因主要是, 国内的镜像加速服务其实是对 `docker hub` 的一次复制, 这里会出现时差, 比如今天 `push` 到 `docker hub` 的镜像, 今天想用, 但国内的镜像加速服务还没更新. 而且库存 <= `docker hub`, 一些多人使用的镜像, `DaoCloud` 才会收录.

现在已经没有选择的余地了......


### 使用 DaoCloud 加速

[Docker Hub Mirror使用手册](http://dockone.io/article/160)

这篇文章的内容可以正常配置 `DaoCloud` 加速, 但是官方说这是 `1.0` 的方案, `2.0` 的方案可以去[官方网站](http://www.daocloud.io/)登录查看.


### DaoCloud 的服务

虽然这次我是使用了 `DaoCloud` 的镜像加速服务, 但其实 `DaoCloud` 还有很多服务, 这里给它一个小广告吧...

- 代码构建
- 镜像仓库, 支持 `DaoCloud` 镜像仓库服务, 包括私人仓库, 企业私有仓库等, 还有 docker hub 的镜像
- 服务集成

