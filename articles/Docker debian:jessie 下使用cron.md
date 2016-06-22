## 07 Dec 15 Docker debian:jessie 镜像使用 cron

-  作者: [sixbyte](http://sixbyte.me/)
-  邮箱: liu.sixbyte@gmail.com
-  禁止转载


服务器上的服务从 `ubuntu` 迁移到 `docker` 上时, 使用的 `debian:jessie` 镜像并没有安装 `cron`, 尝试安装了一下, 遇到了不少问题. `google` 上能找到的关于 `docker` 安装 `cron` 的内容非常少, 迫使我翻阅了大量 `cron` 的资料, 了解到 `cron` 这个工具, 解决了问题, 这里将分享一下.

### 情况

**安装** `debian:jessie` 并没有预装 `cron`

```shell
apt-get update
apt-get install cron
```

**配置当前用户的定时任务**

```shell
crontab -e
crontab -l
# 如下结果
* * * * *   php /var/www/laravel/artisan schedule:run
```

前台模式执行命令

```shell
cron -f
```

理所当然的认为没问题. 然而, 等了N久, 并没有看到命令定时任务 `php /var/www/laravel/artisan schedule:run` 的任何效果.


### 解决办法

```shell
* * * * *   /usr/local/bin/php /var/www/html/autobox/artisan schedule:run
```

### 原因

这是由于环境变量导致的. **`cron` 找不到 `php` 命令**

任务修改如下:

```shell
* * * * *   php /var/www/laravel/artisan schedule:run >> /root/test
```

1分种后 `/root` 目录下有生成 `test` 文件, 但是里面并没有内容, 可以确信, `cron` 是有正常工作的, 问题就在 `crontab` 里面.

翻阅了大量的资料, 感谢这三篇文章带来的提示:

[Linux 下执行定时任务 crontab 命令详解](http://segmentfault.com/a/1190000002628040)
[Linux Crontab定时任务命令详解](http://bbs.csdn.net/topics/390546701)
[crontab 定时任务](http://linuxtools-rst.readthedocs.org/zh_CN/latest/tool/crontab.html)

然后发现了环境变量一说!! 翻阅官方的文档 [crontab.org](http://crontab.org/)

果然!!

```shell
* * * * *   echo $PATH >> /root/test  # 输出 /usr/bin:/bin 而我的 `php` 命令在 `/usr/local/bin/php`
```


### Docker 运行

这里有一篇问题关于 `Docker` 运行 `cron` [Docker: Cronjob is not working](http://stackoverflow.com/questions/24943982/docker-cronjob-is-not-working)

**内容大概如下:**

在容器运行后,可以让容器执行 `cron` 命令运行了, 但容器一旦 `restart` 的时候, `cron` 守护进程被kill, 而且并不会自动生成新的进程运行.

这是因为 Docker is an LXC container, 一个容器只能执行一个命令.需要在前台显式运行 `cron` 守护进程, 要不然就需要 `Supervisor` 或者 `runit`.


然而, `cron` 作为前台运行 `cron -f` 并不像 `php-fpm` 这样的命令有内容输出, 所以, 可以将 `cron` 和 `php-fpm` 组合起来作为一个命令运行. `cron` 命令必须在 `php-fpm` 前.

例子:
```shell
FROM php:5.6-fpm

# cron
RUN apt-get update
RUN apt-get install -y cron
RUN rm -rf /var/lib/apt/lists/*

# myinit.sh
RUN echo '#!/bin/bash' >> /root/myinit.sh
RUN echo "cron && php-fpm" >> /root/myinit.sh
RUN chmod o+x /root/myinit.sh

CMD ["/root/myinit.sh"]
```

