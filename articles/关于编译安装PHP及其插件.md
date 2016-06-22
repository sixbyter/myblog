## 24 Nov 15 关于编译安装PHP及其插件

-  作者: [sixbyte](http://sixbyte.me/)
-  邮箱: liu.sixbyte@gmail.com
-  禁止转载

以前安装过一次, 十分繁琐, 又得翻阅 `php` 的官方文档, 又得google!! 岂有此理


### 编译安装 PHP

**下载**

[官方发布地址](http://php.net/releases/index.php) 选择一个版本下载, 解压.或者:

```shell
wget http://cn2.php.net/distributions/php-5.6.X.tar.gz
```

**安装**

```shell
tar zxvf php-5.6.X.tar.gz
cd php-5.6.X
./configure --prefix=/your/path
make && make install
```
.configure 有一大波参数, 可以通过 ./configure --help 获得, 关于这些参数的解释会发一篇博文

好的, php cli 模式安装完成

### 编译安装 PHP 插件 pcntl 为例

```shell
cd php-5.6.X/etc/pcntl
/your/path/bin/phpize
```
如果出现：

```
Cannot find autoconf. Please check your autoconf installation and the
$PHP_AUTOCONF environment variable. Then, rerun this script.
```
需要先安装 `autoconf`

**MAC 下**
```shell
brew install autoconf
```

**Linux 下, 这个没验证过, 详细的方法 `google` 一下就有**
```shell
wget http://ftp.gnu.org/gnu/autoconf/autoconf-latest.tar.gz
tar zxvf autoconf-latest.tar.gz
cd autoconf-2.69/
./configure –prefix=/data/apps/libs
make &&make install
```


```shell
./configure --with-php-config=/your/path/bin/php-config
make && make install
```

好的, 这时候, 扩展的 `.so` 文件就安装在 `/your/path/lib/php/extensions/no-debug-non-zts-20131226/` 下, 其实有提示.

最后修改 `php.ini`, 这里有个问题, 默认是在 `/usr/local/lib` (/your/path/bin/php --ini 获得), 晕, 我在编译的时候加上 --with-config-file-path=/your/path/lib 都不行!!
最后我修改了 `/usr/local/lib/php.ini`

```shell
extension=/your/path/lib/php/extensions/no-debug-non-zts-20131226/pcntl.so
```

大功告成, 执行:
```shell
/your/path/bin/php -i | grep pcntl 就能看到是否按照了扩展
```