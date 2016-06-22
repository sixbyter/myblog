## 12 Apr 16 HTTPS Let's Encryp 免费的 SSL 证书.md

-  作者: [sixbyte](http://sixbyte.me/)
-  邮箱: liu.sixbyte@gmail.com
-  禁止转载

### 为了你的网站安全

简单来说, 我不想我和访问我网站的用户被运营商或者中间代理 "强奸".

### Let's Encryp

单域名证书其实很便宜, 但是 [Let's Encryp](https://letsencrypt.org/) 免费, 简单, 好用. 更重要的是他的后台强大.

我没有使用 Let's Encrypt 官网提供的工具来申请证书，太大, 太麻烦了. 我更喜欢 [acme-tiny](https://github.com/diafygi/acme-tiny) 这个更为小巧的开源工具.


### 前期准备
1. git, openssl, python. centos 6.5 64 位 python 貌似缺少一些库, 更新: `yum install python-argparse`
2. 确保需要签发的域名使用的DNS为国外DNS服务(如cloudflare,linode等，否则到签发步骤可能会报错)
3. 确保需要签发的域名已指向A记录到你的网站服务器，不要用CNAME等记录，一定要A记录！(否则到签发步骤可能会报错)
4. 最好用ROOT帐号操作(我使用普通帐号操作会报key values mismatch)


### 开始

**我的环境变量**

- 域名: www.sixbyte.me
- 站点的根目录: /usr/local/nginx/html
- acme的地址: /usr/local/nginx/html/acme-tiny


nginx 的配置

```bash
# https 配置
#server {
#    listen       443 ssl;
#    server_name  www.sixbyte.me;

#    ssl_certificate      /usr/local/nginx/html/acme-tiny/chained.pem;
#    ssl_certificate_key  /usr/local/nginx/html/acme-tiny/www.sixbyte.me.key;

#    ssl_session_cache    shared:SSL:1m;
#    ssl_session_timeout  5m;
#
#    ssl_ciphers  HIGH:!aNULL:!MD5;
#    ssl_prefer_server_ciphers  on;
#
#    location / {
#        root   html;
#        index  index.html index.htm;
#    }
#}

server {
    listen 80;
    server_name www.sixbyte.me sixbyte.me;

    location / {
        rewrite ^/(.*)$ https://www.sixbyte.me/$1 permanent;
    }

    location ^~ /.well-known/acme-challenge/ {
        alias /usr/local/nginx/html/acme-tiny/;
        try_files $uri =404;
    }
}
```

执行以下命名

```bash
cd /usr/local/nginx/html

git clone https://github.com/diafygi/acme-tiny.git

cd acme-tiny

# 用户私钥: account.key
openssl genrsa 4096 > account.key

# 域名私钥: www.sixbyte.me.key
openssl genrsa 4096 > www.sixbyte.me.key

# 生成 www.sixbyte.me.crt
openssl req -new -sha256 -key www.sixbyte.me.key -subj "/CN=www.sixbyte.me" > www.sixbyte.me.csr
# 多域名如下:
# openssl req -new -sha256 -key www.sixbyte.me.key -subj "/" -reqexts SAN -config <(cat /etc/ssl/openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:sixbyte.me,DNS:www.sixbyte.me")) > www.sixbyte.me.csr
# 如果报错找不到 openssl.cnf, 别慌, 执行如下: wget http://web.mit.edu/crypto/openssl.cnf > /etc/ssl/openssl.cnf

# 申请证书, 在Let’s Encrypt 服务器提交签发证书(程序大致操作：acme-tiny.py会生成一个密钥文件到acme-tiny目录下，然后Let’s Encrypt 证书签发服务器会访问需签发域名/.well-known/acme-challenge/路径下acme-tiny.py生成的密钥文件), 所以DNS解析在国内可能会报错.
python acme_tiny.py --account-key account.key --csr www.sixbyte.me.csr --acme-dir /usr/local/nginx/html/acme-tiny/ > www.sixbyte.me.crt

# 合并 Let's Encrypt 的中间证书
wget -O - https://letsencrypt.org/certs/lets-encrypt-x1-cross-signed.pem > intermediate.pem
cat www.sixbyte.me.crt intermediate.pem > chained.pem
```

取消上面nginx 的 https 注释, 然后测试:
```bash
# /usr/local/nginx/sbin/nginx -t
nginx -t
nginx: the configuration file /usr/local/nginx/conf/nginx.conf syntax is ok
nginx: configuration file /usr/local/nginx/conf/nginx.conf test is successful
```

通过后重启nginx
```bash
service nginx restart
```


最后记得注释这段:
```bash
server {
    listen 80;
    server_name www.sixbyte.me sixbyte.me;

    location / {
        rewrite ^/(.*)$ https://www.sixbyte.me/$1 permanent;
    }

    #location ^~ /.well-known/acme-challenge/ {
    #    alias /usr/local/nginx/html/acme-tiny/;
    #    try_files $uri =404;
    #}
}
```