## 29 Jan 16 基于Docker的PHP开发环境 PHPWEB

-  作者: [sixbyte](http://sixbyte.me/)
-  邮箱: liu.sixbyte@gmail.com
-  禁止转载


开发环境每个人都有一套, 我自己也有一套. 一开始学习 `PHP` 的时候用 `WAMP`. 类似的有很多, 初学者能快速的上手. 多亏了这些集成环境, 降低了学习的难度. 但是一旦在工作了, 要求 `PHP` 的生产环境一般都在 `Linux` 系统上, 开发为了环境统一, 都尽量在 `Linux` 上开发. 这时, `LNMP` 这样的集成环境也是一个不错的选择.

但是呢, 我一直无法理解 `apt-get install php5` 这个命令, 因为我不知道它替我做了多少事情?( `php.ini` 的默认目录, `php-fpm.conf` 的默认目录, 等等). 有洁癖的我, 就自己折腾了一番. `Docker` 也是我学习的一个领域, 而且非常喜欢, 解决了我洁癖的问题, 安装过程还记录在 `Dockerfile` 里, 非常适合我. 自己折腾了几套基于 `Docker` 的 `PHP` 的开发环境, 最近整理折腾的这套最喜欢, 介绍一下.

[PHPWEB 项目地址](https://github.com/sixbyter/phpweb)

由于项目的更新会在 `github` 上, 而且可能会比较频繁, 确保时效性, 所以这里不做详细介绍.

### 简介

![结构](https://raw.githubusercontent.com/sixbyter/phpweb/master/doc/F035DBB5-CE77-4C0C-8829-542A6C4F9AEE.png)

PHPWEB 安装了以下内容:

- php 7.0.2
- composer
- laravel/envoy
- laravel/installer
- nodejs 5.5.0
- npm 3.3.12
- nginx 1.9.9
- gulp
- bower
- grunt-cli
- pm2
- supervisor
- cron
- git


这套环境非常适合 `Laravel` 框架开发, 同时跟 `Homestead` 非常相似. 其他的 `Homestead` 服务只需要运行相应的服务的容器, 然后链接就能提供服务.

如果你要部署到生产环境, 为了运行速度你可以修改 `Dockerfile`, 删减一些没有必要的安装, 重新 `docker build` 一个镜像. 这个环境的构建是透明的, 随意更改, 切换软件版本等.


### 特性

因为基于 `docker`搭建的镜像, "运输"非常的方便. 开发 `sixbyte/hello-world` 为例

开发时构建环境:

```
docker-compose up -d
```

开发过程中安装的依赖, 配置全部记录下来, 完成后创建项目的 `Dockerfile`

```
FROM sixbyte/phpweb

# 安装依赖

# 下载代码
RUN git clone xxxxxxxxxxx

# 配置
RUN cp xxxxxxx /etc/nginx/conf.d/phpweb.conf

# .....
```

我们尽可能在构建镜像做更多的事情, 让运行镜像更加精简, 比如运行只需要 `docker run -d sixbyte/hello-world:0.1`.

构建镜像, 其中 0.1 我们称为是项目的版本号:

```
docker build -t sixbyte/hello-world:0.1 .
```

这样, 我们 `0.1` 版本的项目就是镜像 `sixbyte/hello-world:0.1`, `push` 或者 `dockerfile` 的方式送到测试环境和生产环境. 这样开发, 测试, 生产的运行环境是一致的, 而且实现了对项目的版本控制, 的更新和回滚都是非常方便的.