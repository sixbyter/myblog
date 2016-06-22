## 22 Dec 15 PHP 扩展库学习进度

-  作者: [sixbyte](http://sixbyte.me/)
-  邮箱: liu.sixbyte@gmail.com
-  禁止转载

参考地址: [PHP 扩展 按归属分](http://php.net/manual/zh/extensions.membership.php)

### 核心扩展库

以下并非真正的扩展库，它们属于 PHP 内核的一部分，不能通过编译选项将其排除。

| 扩展               | 阅读次数 | 掌握程度 |
| ----------------- | ------- | ------- |
| 数组               |    2    |    长     |
| 类/对象            |    1    |         |
| 日期/时间           |   1     |         |
| 目录               |    1    |         |
| 错误处理            |    1    |         |
| 程序执行            |    1    |         |
| Filesystem        |    1    |      长   |
| Filter            |    1    |         |
| Function Handling |    1    |         |
| Hash              |    1    |         |
| PHP 选项/信息       |   1     |    长     |
| Mail              |         |         |
| Math              |         |     长    |
| Misc.             |     1   |         |
| 网络               |     1   |    长    |
| 输出控制            |     1   |         |
| Password Hashing  |     1   |   奇怪呀      |
| Phar              |         |    长     |
| 反射               |         |    长     |
| POSIX Regex       |    0    | 不推荐使用, php7移除该函数 |
| Sessions          |     1   |         |
| SPL               |     1   |    长     |
| Streams           |     1   |      长   |
| 字符串             |     1   |     长..奇葩的php..看到想吐..超多函数, 硬生生看完, 写完demo, 有些函数很蛋疼..    |
| Tokenizer         |     1   |         |
| URLs              |    1    |         |
| Variable handling |    1    |    长     |

可恶, fflush 和 stream_set_write_buffer 如何能写出能看出效果的demo!!?



### 绑定的扩展库

以下扩展库绑定在 PHP 发行包中。

| 扩展               | 阅读次数 | 掌握程度 |
| ----------------- | ------- | ------- |
| Apache            |         |         |
| BC Math           |         |         |
| Calendar          |         |         |
| COM               |         |         |
| Ctype             |         |         |
| DBA               |         |         |
| Exif              |         |         |
| Fileinfo          |         |         |
| FTP               |         |         |
| iconv             |         |         |
| GD                |         |         |
| intl              |         |         |
| JSON              |    1    |         |
| 多字节字符串         |         |         |
| NSAPI             |         |         |
| PCNTL             |         |         |
| PCRE              |         |         |
| PDO               |         |         |
| POSIX             |         |         |
| Semaphore         |         |         |
| Shared Memory     |         |         |
| Sockets           |    1    |         |
| SQLite3           |         |         |
| XML-RPC           |         |         |
| Zlib              |         |         |


### 外部扩展库

这些扩展库已经绑定在 PHP 发行包中，但是要编译以下扩展库，需要外部的库文件。

| 扩展                   | 阅读次数 | 掌握程度  |
| --------------------- | ------- | ------- |
| Bzip2                 |         |         |
| cURL                  |    1    |         |
| dBase                 |         |         |
| DOM                   |         |         |
| Enchant               |         |         |
| FrontBase             |         |         |
| Gettext               |         |         |
| GMP                   |         |         |
| Firebird/InterBase    |         |         |
| Informix              |         |         |
| IMAP                  |         |         |
| LDAP                  |         |         |
| libxml                |         |         |
| Mcrypt                |         |         |
| Mhash                 |         |         |
| mSQL                  |         |         |
| Mssql                 |         |         |
| Mysql                 |         |         |
| Mysqli                |         |         |
| Mysqlnd               |         |         |
| OCI8                  |         |         |
| OpenSSL               |         |         |
| MS SQL Server (PDO)   |         |         |
| Firebird (PDO)        |         |         |
| MySQL (PDO)           |         |         |
| Oracle (PDO)          |         |         |
| ODBC and DB2 (PDO)    |         |         |
| PostgreSQL (PDO)      |         |         |
| SQLite (PDO)          |         |         |
| PostgreSQL            |         |         |
| Pspell                |         |         |
| Readline              |         |         |
| Recode                |         |         |
| SimpleXML             |         |         |
| SNMP                  |         |         |
| SOAP                  |         |         |
| Sybase                |         |         |
| Tidy                  |         |         |
| ODBC                  |         |         |
| WDDX                  |         |         |
| XML 解析器             |         |         |
| XMLReader             |         |         |
| XMLWriter             |         |         |
| XSL                   |         |         |
| Zip                   |         |         |


### PECL 扩展库

这里不列出
