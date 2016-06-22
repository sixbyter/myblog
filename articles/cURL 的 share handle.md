## 25 Nov 15 PHP cURL 的 share handle

-  作者: [sixbyte](http://sixbyte.me/)
-  邮箱: liu.sixbyte@gmail.com
-  禁止转载


最近扫了一波 [鸟哥](http://www.laruence.com/) 的博客里我现阶段能看得懂的文章, 收益匪浅, 收藏和研究了一些感兴趣的列子. 讲到 `php` 更新了 `cURL` 的扩展库 `libcurl`, 于是去官网看了官网文档, 发现多了一个 curl_share_init 的函数, 第一次见, 看了介绍后相当兴奋, 让我想起了以前cURL做模拟登录的例子

### 配置
`curl_share_setopt` 现在只支持 `CURLSHOPT_SHARE` 和  `CURLSHOPT_UNSHARE` 两种选项

**选项**

Option              | Description
--------------------|--------------
CURLSHOPT_SHARE     | 指定一个会被分享的数据类型
CURLSHOPT_UNSHARE   | 指定一个不被分享的数据类型

**值**

Value                       | Description
----------------------------|--------------
CURL_LOCK_DATA_COOKIE       | 分享 cookie
CURL_LOCK_DATA_DNS          | 分享 dns 缓存, 当使用 cURL multi handles 时, 默认所有连接都会使用 dns 缓存
CURL_LOCK_DATA_SSL_SESSION  | 分享 SSL session IDs, 当重新访问同一个服务器时, 减少和 SSL 握手的消耗

其中分享 dns 缓存 和 分享 SSL session IDs 都很好理解, 都是提升性能的, 这里不做演示. 我想说的是 `CURL_LOCK_DATA_COOKIE`, 不但提升性能, 还能写出更优雅的代码.

**以前的模拟登录**
```php
<?php
/**
* Curl 模拟登录 discuz 程序
* 尚未实现开启验证码的的论坛登录功能
*/

$discuz_url = 'http://xxxxx.discuz.com';//论坛地址
$login_url = $discuz_url .'/logging.php?action=login';//登录页地址
$get_url = $discuz_url .'/my.php?item=threads'; //我的帖子

// 编辑post 数据, 通过抓包获取
$post_fields = [];
$post_fields['loginfield'] = 'username';
$post_fields['loginsubmit'] = 'true';
$post_fields['username'] = 'xxxxx';
$post_fields['password'] = 'xxxxx';
$post_fields['questionid'] = 0;
$post_fields['answer'] = '';
$post_fields['seccodeverify'] = '';

// 创建一个临时文件, 记录登录的 cookie
$cookie_file = dirname(__FILE__) . '/cookie.txt';
$ch = curl_init($login_url);
curl_setopt($ch, CURLOPT_HEADER, 0);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
curl_setopt($ch, CURLOPT_POST, 1);
curl_setopt($ch, CURLOPT_POSTFIELDS, $post_fields);
curl_setopt($ch, CURLOPT_COOKIEJAR, $cookie_file);
curl_exec($ch);
curl_close($ch);

//带着上面得到的COOKIE获取需要登录后才能查看的页面内容
$ch = curl_init($get_url);
curl_setopt($ch, CURLOPT_HEADER, 0);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 0);
curl_setopt($ch, CURLOPT_COOKIEFILE, $cookie_file);
$contents = curl_exec($ch);
curl_close($ch);
var_dump($contents);
?>
```

ok.... 非常麻烦, 要创建一个文件保存 `cookie`, 但有了 `CURL_LOCK_DATA_COOKIE`, 我们可以这样

```php
<?php

$sh = curl_share_init();
curl_share_setopt($sh, CURLSHOPT_SHARE, CURL_LOCK_DATA_COOKIE);

$ch1                        = curl_init("http://sixbyte.sinaapp.com/wp-login.php");
$post_fields                = [];
$post_fields['log']         = "不告诉你";
$post_fields['pwd']         = "不告诉你";
$post_fields['redirect_to'] = "http://sixbyte.sinaapp.com/wp-admin/";
$post_fields['testcookie']  = 1;
curl_setopt($ch1, CURLOPT_HEADER, 0);
curl_setopt($ch1, CURLOPT_RETURNTRANSFER, 1);
curl_setopt($ch1, CURLOPT_POST, 1);
curl_setopt($ch1, CURLOPT_POSTFIELDS, $post_fields);
curl_setopt($ch1, CURLOPT_SHARE, $sh);

curl_exec($ch1);

$ch2 = curl_init("http://sixbyte.sinaapp.com/wp-admin/");
curl_setopt($ch2, CURLOPT_SHARE, $sh);

echo curl_exec($ch2);

curl_share_close($sh);

curl_close($ch1);
curl_close($ch2);

```

是不是觉得超简单的!