## 14 Feb 16 GO DNS协议及本地DNS服务器(一)

-  作者: [sixbyte](http://sixbyte.me/)
-  邮箱: liu.sixbyte@gmail.com
-  禁止转载

那天老胡给我推荐了一个`dns`工具[Pcap_DNSProxy](https://github.com/chengr28/Pcap_DNSProxy), 据说可以防止 `dns` 污染. 然后很好奇, 什么是 `dns` 污染, 我可是只知道 `dns` 劫持啊(其实`dns` 劫持是`dns`污染的一种). 还有这软件是怎样做到, 乱搞了一堆, 不知不觉去触碰到了 `dns` 协议.

### 测试 `DNS` 服务器

测试默认的 `dns` 服务器
```
dig baidu.com
```
测试指定的 `dns` 服务器
```
dig baidu.com @127.0.0.1
```

### DNS 客户端解析过程

以下都是我的理解, 官方的说法请参考 `RFC` 文档关于 `DNS` 的章节, 比如 1034, 1035.

只解析 `DNS` 的客户端解析过程, 因为 `DNS` 解析大家都很清楚, 是一个递归过程, 本地查询不到记录会访问根(上级) `DNS` 服务器.

Chrome 打开 `google.com` 为例.

-> google.com
-> 代理
-> Chorme 的 `dns` 缓存
-> `/etc/hosts`
-> 系统DNS缓存
-> `/etc/resolv.conf` 的 `DNS` 服务器

当使用 `shadowsocks` 作代理时, 开启全局代理模式时, `/etc/hosts` 文件的映射记录 not work.
修改 `/etc/hosts` 后, 浏览器为什么没有实时更新. chrome 可以访问 `chrome://net-internals/#dns` 清理浏览器的缓存.
系统DNS缓存是部分系统才有的, 苹果官方提供了一个[清理 `OS X` 的命令](https://support.apple.com/zh-cn/HT202516), 在 OS X v10.10.4 或更高版本中 `sudo killall -HUP mDNSResponder` 但是啊~~我自己测试的时候并没有发现系统有 `dns` 缓存, 很奇怪.
`/etc/resolv.conf` 这个可以自行修改. 没有定义默认记录的是所连网络的地方的运营商的 `DNS` 服务器. 比如我在公司连的是深圳电信的 `DNS` 服务器.

```shell
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 46921
;; flags: qr rd ra; QUERY: 1, ANSWER: 0, AUTHORITY: 1, ADDITIONAL: 0
#
# Mac OS X Notice
#
# This file is not used by the host name and address resolution
# or the DNS query routing mechanisms used by most processes on
# this Mac OS X system.
#
# This file is automatically generated.
#
nameserver 202.96.134.133
nameserver 202.96.128.86
```

### DNS 协议

参考资料:
- [DNS协议Golang实现](http://blog.cyeam.com/network/2015/02/03/dns/)
- [DNS协议分析](http://blog.cyeam.com/network/2015/01/29/dns/)
- [RFC1035-Domain Implementation and Specification](https://www.rfc-editor.org/rfc/rfc1035.txt)

写不下去了!!! RFC1035 这么长...最好还是去看看, 实在很难用三言两语去描述...

Domain System 比较官方的描述是 `RFC` 文档的两个姐妹文档1034和1035, 而根据互联网的发展, 后面还有十几篇补充文档. 而 1035 是最详细的描述 DNS 协议和实现的. 英文不太好的可以尝试参考我这里关于 `DNS` 协议的解释, 但不代表完全准确, 我也是阅读了1035原文和一些翻译以及 `DNS` 的国人文章的一些总结.

- 基于 `TCP` 或 `UDP`

DNS 协议没有强硬要求使用某种传输层协议. 但目前一般使用 `UDP` 协议. 端口号 `53`

- RR 定义

不要问我 `RR` 是什么, 我也不知道全称...要是你知道请告诉我....

下面的格式根据16位的格式显示, 其实 `//` 表示不定长

```
                                    1  1  1  1  1  1
      0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                                               |
    /                                               /
    /                      NAME                     /
    |                                               |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                      TYPE                     |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                     CLASS                     |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                      TTL                      |
    |                                               |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                   RDLENGTH                    |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
    /                     RDATA                     /
    /                                               /
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
```

NAME: 域名, 比如 sixbyte.me
TTL: Time To Live, 生存时间, 也就是域名的过期时间
RDLENGTH: RDATA 的长度
RDATA: 内容主体

**TYPE**

| TYPE    |  value  | meaning                                    |
|---------|---------|--------------------------------------------|
| A       |   1     | a host address                             |
| NS      |   2     | an authoritative name server               |
| MD      |   3     | a mail destination (Obsolete - use MX)     |
| MF      |   4     | a mail forwarder (Obsolete - use MX)       |
| CNAME   |   5     | the canonical name for an alias            |
| SOA     |   6     | marks the start of a zone of authority     |
| MB      |   7     | a mailbox domain name (EXPERIMENTAL)       |
| MG      |   8     | a mail group member (EXPERIMENTAL)         |
| MR      |   9     | a mail rename domain name (EXPERIMENTAL)   |
| NULL    |   10    | a null RR (EXPERIMENTAL)                   |
| WKS     |   11    | a well known service description           |
| PTR     |   12    | a domain name pointer                      |
| HINFO   |   13    | host information                           |
| MINFO   |   14    | mailbox or mail list information           |
| MX      |   15    | mail exchange                              |
| TXT     |   16    | text strings                               |
| AAAA    |         |                                            |
| UINFO   |         |                                            |
| UID     |         |                                            |
| GID     |         |                                            |
| ANY     |         |                                            |

**CLASS**

| CLASS   |  value  | meaning                                    |
|---------|---------|--------------------------------------------|
| IN      |   1     | the Internet                               |
| CS      |   2     | the CSNET class (Obsolete - used only for examples in some obsolete RFCs) |
| CH      |   3     | the CHAOS class                            |
| HS      |   4     | Hesiod [Dyer 87]                           |

响应有一个 192 16 或者 c0 0c 的2字节, 其实是用来表示域名, 跟偏移量有关, 详细: [对DNS报文的理解](http://jianjian.blog.51cto.com/blog/35031/5170)


### 本地 DNS 服务器

为什么要开发一个本地的 `DNS` 的本地服务器?!!? 还用问嘛??! 当然是因为对现在的 `DNS` 服务器很不满啊!! 两方面, 一个是速度, 一个是安全.

刚好接触到一个 `GO` 写的 `DNS` 本地服务器 [godns](https://github.com/kenshinx/godns) 听讲很厉害呢, 能做 `dns` 本地缓存, 还能指定映射.
然后我马上试了一下, 卧槽!!! 没有使用 `godns` 的时候, `dig baidu.com` 平均要8ms, 而且经常超时到1000ms; 使用后...0ms. 这甜头, 我怎么不去咬一口.

然后老胡介绍我的[Pcap_DNSProxy](https://github.com/chengr28/Pcap_DNSProxy)也很在意, 纠结是怎样做到防污染的, 毕竟 `DNS` 污染是在可怕, 那些运营商的 `dns` 服务器是在不可靠(都怪你们劫持投放广告!!). 还有那个 `chengr28` 好像是个高中生来着??! 吓我一跳.

因为这些, 我有了开发一个自己的 `dns` 的服务器, 不是 `godns` 不好, 就是好奇了, 想折腾一下, 而且不是自己写的总不放心. 然后你就得去了解 `dns` 协议


// []byte 函数并没有引用, 函数里可以修改其值

func (this *dnsMsg) Pack() []byte {
    bs := make([]byte, 12)

    binary.BigEndian.PutUint16(bs[0:2], this.ID)

    fmt.Println(reflect.TypeOf(bs))

    var a [2]uint8
    b := new([2]uint8)
    fmt.Println(reflect.TypeOf(a))
    fmt.Println(reflect.TypeOf(a[0:2]))
    fmt.Println(reflect.TypeOf(a))
    fmt.Println(reflect.TypeOf(b))
    fmt.Println(reflect.TypeOf(b[0:2]))
    test(a[0:2])

    var s [2]string
    test2(s[0:2])
    fmt.Println(s)

    fmt.Println(a)
    return bs
}

func test(a []uint8) {
    a[0] = 1
    a[1] = 1
}

func test2(a []string) {
    a[0] = "1"
    a[1] = "2"
}

func (this *dnsMsg) Pack() []byte {
    bs := make([]byte, 12)
    binary.BigEndian.PutUint16(bs[0:2], this.ID)

    a := []int{1, 2, 3}
    fmt.Println(a)
    modifySlice(a)
    fmt.Println(a)
    fmt.Println(reflect.TypeOf(a))

    return bs
}

func modifySlice(data []int) {
    data = nil
}


不同的DNS服务器返回的结果不一样
抽风的114.114.114.114

www.baidu.com
www 是 name
baidu com 都是label
baidu 是二级域名, 是这个域名的主体
com 是后缀, 是顶级域名


labels          63 octets or less

names           255 octets or less

TTL             positive values of a signed 32 bit number.

UDP messages    512 octets or less


RRs

QTYPE 是 TYPE 的一个超集,
AXFR            252 A request for a transfer of an entire zone

MAILB           253 A request for mailbox-related records (MB, MG or MR)

MAILA           254 A request for mail agent RRs (Obsolete - see MX)

*               255 A request for all records


[IN-ADDR.ARPA](http://blog.sina.com.cn/s/blog_537d395001000byz.html) 具体用go写一个demo才行


   - a sequence of labels ending in a zero octet

   - a pointer

   - a sequence of labels ending with a pointer