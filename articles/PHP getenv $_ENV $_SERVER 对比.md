## 15 Jan 16 PHP getenv $_ENV $_SERVER 对比

-  作者: [sixbyte](http://sixbyte.me/)
-  邮箱: liu.sixbyte@gmail.com
-  禁止转载


翻阅 `php` 的官方文档, 对 `php` 中的 `getenv` `$_ENV` `$_SERVER` 产生了好奇. 因为他们内容十分相似, 但却独立存在. 同时 `Laravel` 在设置环境变量的时候全都设置一遍, 让我十分疑惑. 这里让我们对比一下他们的区别.

### 首先
我们都清楚, 在不同的环境下, 它们的情况会不同. 特别是在不同的 `sapi` 模式下. 本文章只对 `cli` 模式和 `fpm-fcgi` 模式进行分析.
`php.ini` 中, 注意 `variables_order` 的设置为默认的 `EGPCS`

### 值的对比

首先, 获取所有的 `getenv` 键值对. 官方文档说得很清楚, 所有的 `getenv` 设置都在 `phpinfo()` 里.
```php
function get_all_info_env()
{
    ob_start();
    phpinfo(INFO_ENVIRONMENT);
    $phpinfo_env_content = ob_get_clean();
    $info_env            = [];
    $sapi                = php_sapi_name();
    if ($sapi === 'cli') {
        $lines = explode("\n", $phpinfo_env_content);
        foreach ($lines as $line) {
            if (strpos($line, "=>") === false) {
                continue;
            }
            $e                     = explode("=>", $line);
            $info_env[trim($e[0])] = trim($e[1]);
        }
        unset($info_env['Variable']);
    }
    if ($sapi === 'fpm-fcgi' || $sapi === 'cli-server') {
        preg_match_all('#<tr><td class="e">(.*)</td><td class="v">(.*)</td></tr>#', $phpinfo_env_content, $out, PREG_SET_ORDER);
        foreach ($out as $line) {
            $info_env[$line[1]] = $line[2];
        }

    }
    return $info_env;
}
```

然后对比 `$_ENV`, `$_SERVER`, `get_all_info_env()` 三个数组, 可以看到如下结果:

**cli 模式**

- `$_ENV` === `get_all_info_env()` 都是 `/bin/bash` 等外壳程序的环境变量.
- `$_SERVER` 是 `$_ENV` 的超集. `$_SERVER` 多出的部分如下. 我的脚本名为 `1.php`

```
  ["PHP_SELF"]=>
  string(5) "1.php"
  ["SCRIPT_NAME"]=>
  string(5) "1.php"
  ["SCRIPT_FILENAME"]=>
  string(5) "1.php"
  ["PATH_TRANSLATED"]=>
  string(5) "1.php"
  ["DOCUMENT_ROOT"]=>
  string(0) ""
  ["REQUEST_TIME_FLOAT"]=>
  float(1452838923.5402)
  ["REQUEST_TIME"]=>
  int(1452838923)
  ["argv"]=>
  array(1) {
    [0]=>
    string(5) "1.php"
  }
  ["argc"]=>
  int(1)
```

**cli-server 模式**

- 同样, `$_ENV` === `get_all_info_env()` 都是 `/bin/bash` 等外壳程序的环境变量.
- `$_SERVER` 和 `$_ENV` 交集为空.

`$_SERVER` 更多是和请求有关, 可以被 `webserver` 修改, `PATH_INFO`等



**fpm-fcgi 模式**

- 总体来说, `$_ENV` 是 `get_all_info_env()` 的超集, `$_SERVER` 是 `$_ENV` 的超集
- `$_ENV`, `$_SERVER` 还带有 `web server` 带来的信息, 比如 `nginx` 的 `SCRIPT_FILENAME` 等

`php` 的官方文档也有这方面的解释. [地址](http://us.php.net/manual/zh/ini.core.php#ini.variables-order)
> In both the CGI and FastCGI SAPIs, $_SERVER is also populated by values from the environment; S is always equivalent to ES regardless of the placement of E elsewhere in this directive.

但是, 一个一个奇怪的现象(只出现在我的fpm镜像, 我使用的是docker官网的fpm镜像), 一旦 `$_SERVER` 在代码里比$_ENV早出现. `$_ENV` 和 `$_SERVER` 两个数组完全相等.

```
$_SERVER;
$_ENV;
var_dump($_ENV === $_SERVER); // true
var_dump(array_diff_assoc($_SERVER, $_ENV)); // array() 空数组
```

### 设置值

`$_ENV` 和 `getenv` 是独立存在的. 修改 `$_ENV` 不会影响 `getenv` 的值. 修改 `getenv` 不会影响 `$_ENV` 的值.
`$_SERVER` 当然也是独立的了.




整体来说, `fpm-fcgi` 模式比较混乱, 编写框架程序, 服务器设计时需要注意的地方很多.
