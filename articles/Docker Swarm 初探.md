## 30 Jun 16 Docker Swarm 初探

-  作者: [sixbyte](http://sixbyte.me/)
-  邮箱: liu.sixbyte@gmail.com
-  禁止转载


`docker 1.12 RC2` 出来了, 很多新功能已经可以试用了. 其中 `docker swarm` 吸引了我注意. 因为我最近也有搭建过 `swarm` 集群. 但是 `docker1.12` 在 `swarm` 上有很大的变动. `swarm` 成为了 `docker` 的子命令. 这是一个历史性的时刻.

> 官网文档.(虽然还不完善, 但是介绍得很好.)
> https://docs.docker.com/engine/swarm/

### Docker Swarm 特性

#### 来自官方的特性介绍
 
> - Cluster management integrated with Docker Engine (集群的编排和管理集成到 Docker Engine)
> - Decentralized design (分散设计, 这个我没看懂啊..我粗糙地翻译一下. 不是在部署时区分不同的节点角色, 而是Docker在运行时处理任何专门化..大雾. 你可以用 Docker > 部署2种节点角色, managers和workers, 意思是你可以在同一个单一的磁盘镜像创建完整的集群....大雾..)
> - Declarative service model(定义服务模型..终于来了, 集群里的重要概念)
> - Scaling (规模/等级, 可以为服务设定规模. 也就是一个服务设置规模是5, 那么会创建5个容器为这个服务服务.)
> - Desired state reconciliation (期望状态调和, 大雾..某服务的规模是5, 如果有一个容器挂了, swarm会自动开启一个新的容器顶上. 节点挂了, > 该节点上要跑的服务会转移到别的节点上.)
> - Multi-host networking (多主机容器网络)
> - Service discovery (服务发现, 这个再也不用 Consul/Etcd 等第三方的支持了, Docker Engine 自带)
> - Load balancing (负载均衡)
> - Secure by default (默认安全)
> - Rolling updates (滚动更新)

#### 新增的子命令

docker swarm

```
$ docker swarm --help

Usage:  docker swarm COMMAND

Manage Docker Swarm

Options:
      --help   Print usage

Commands:
  init        Initialize a Swarm
  join        Join a Swarm as a node and/or manager
  update      Update the Swarm
  leave       Leave a Swarm
  inspect     Inspect the Swarm

Run 'docker swarm COMMAND --help' for more information on a command.
```

docker node

```
$ docker node --help

Usage:  docker node COMMAND

Manage Docker Swarm nodes

Options:
      --help   Print usage

Commands:
  accept      Accept a node in the swarm
  demote      Demote a node from manager in the swarm
  inspect     Inspect a node in the swarm
  ls          List nodes in the swarm
  promote     Promote a node to a manager in the swarm
  rm          Remove a node from the swarm
  tasks       List tasks running on a node
  update      Update a node

Run 'docker node COMMAND --help' for more information on a command.
```

docker service

```
$ docker service --help

Usage:  docker service COMMAND

Manage Docker services

Options:
      --help   Print usage

Commands:
  create      Create a new service
  inspect     Inspect a service
  tasks       List the tasks of a service
  ls          List services
  rm          Remove a service
  scale       Scale one or multiple services
  update      Update a service

Run 'docker service COMMAND --help' for more information on a command.
```

docker stack (实验性)

```
$ docker stack --help

Usage:  docker stack COMMAND

Manage Docker stacks

Options:
      --help   Print usage

Commands:
  config      Print the stack configuration
  deploy      Create and update a stack
  rm          Remove the stack
  tasks       List the tasks in the stack

Run 'docker stack COMMAND --help' for more information on a command.

让我们来看看这些特性.
```

docker deploy (实验性)

```
$ docker deploy --help

Usage:  docker deploy [OPTIONS] STACK

Create and update a stack

Options:
  -f, --bundle string   Path to a bundle (Default: STACK.dsb)
      --help            Print usage
```


### 简单试用-搭建一个集群

效果是这样的, 3个节点, 其中1个 `manager`, 2个 `worker`.

创建机器. 其中 `--engine-opt="registry-mirror=http://xxxxxx.m.daocloud.io"` 是我 `DaoCloud` 镜像仓库加速器的地址. 主要是为了解决国内镜像下载速度慢的问题.

```sh
docker-machine create  -d virtualbox --engine-opt="registry-mirror=http://xxxxxx.m.daocloud.io" manager
docker-machine create  -d virtualbox --engine-opt="registry-mirror=http://xxxxxx.m.daocloud.io" worker1
docker-machine create  -d virtualbox --engine-opt="registry-mirror=http://xxxxxx.m.daocloud.io" worker2
```

```sh
$ docker-machine ls
NAME      ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER        ERRORS
manager   *        virtualbox   Running   tcp://192.168.99.100:2376           v1.12.0-rc2
worker1     -        virtualbox   Running   tcp://192.168.99.101:2376           v1.12.0-rc2
worker2     -        virtualbox   Running   tcp://192.168.99.102:2376           v1.12.0-rc2
```

创建集群, 初始化集群. 初始化之前我们打印一下 netstat 看看创建集群的时候会多了哪些网络连接.

```sh
$ docker-machine ssh manager

docker@manager:~$ netstat -an

```
初始化集群后, 我们可以看到除了会监听 2377 端口外, 还有监听 7946 端口, 还有 4789 端口. 同时还会本机和 2377 端口建立连接.....

> TCP port 2377 for cluster management communications

> TCP and UDP port 7946 for communication among nodes

> TCP and UDP port 4789 for overlay network traffic

```sh
docker@manager:~$ docker swarm init --listen-addr 192.168.99.100:2377
Swarm initialized: current node (cmlouhnpzmmlb59526loze8v3) is now a manager.

docker@manager:~$ netstat -an
Active Internet connections (servers and established)
Proto Recv-Q Send-Q Local Address           Foreign Address         State
tcp        0      0 192.168.99.100:2377     0.0.0.0:*               LISTEN
tcp        0      0 192.168.99.100:7946     0.0.0.0:*               LISTEN
tcp        0      0 192.168.99.100:40568    192.168.99.100:2377     ESTABLISHED
tcp        0      0 192.168.99.100:2377     192.168.99.100:40568    ESTABLISHED
udp        0      0 0.0.0.0:4789            0.0.0.0:*
udp        0      0 192.168.99.100:7946     0.0.0.0:*
```

另外2台节点加入集群.

```sh
$ docker-machine ssh worker1
docker@worker1:~$ docker swarm join 192.168.99.100:2377
This node joined a Swarm as a worker.

$ docker-machine ssh worker2
docker@worker2:~$ docker swarm join 192.168.99.100:2377
This node joined a Swarm as a worker.
```
我们执行这些命令查看一下集群的信息.

```sh
$ docker-machine ssh manager
```

```sh
docker@manager:~$ docker info
Containers: 0
 Running: 0
 Paused: 0
 Stopped: 0
Images: 0
Server Version: 1.12.0-rc2
Storage Driver: aufs
 Root Dir: /mnt/sda1/var/lib/docker/aufs
 Backing Filesystem: extfs
 Dirs: 0
 Dirperm1 Supported: true
Logging Driver: json-file
Cgroup Driver: cgroupfs
Plugins:
 Volume: local
 Network: null bridge overlay host
Swarm: active
 NodeID: cmlouhnpzmmlb59526loze8v3
 IsManager: Yes
 Managers: 1
 Nodes: 3
 CACertHash: sha256:b0d128fef0cae30bc593ff60c8b935615d5994c2d6c2d387836ad5ebeeaa7f92
Runtimes: default
Default Runtime: default
Security Options: seccomp
Kernel Version: 4.4.13-boot2docker
Operating System: Boot2Docker 1.12.0-rc2 (TCL 7.1); HEAD : 52952ef - Fri Jun 17 21:01:09 UTC 2016
OSType: linux
Architecture: x86_64
CPUs: 1
Total Memory: 995.9 MiB
Name: manager
ID: 4C6Y:AG2E:HE6U:XYJB:L5WK:4E7O:JJ4G:5ZZZ:5WUW:A5GT:K57A:WZUK
Docker Root Dir: /mnt/sda1/var/lib/docker
Debug Mode (client): false
Debug Mode (server): true
 File Descriptors: 37
 Goroutines: 125
 System Time: 2016-06-29T13:25:10.488422335Z
 EventsListeners: 0
Registry: https://index.docker.io/v1/
Labels:
 provider=virtualbox
Insecure Registries:
 127.0.0.0/8
```

```sh
docker@manager:~$ docker node ls
ID                           NAME     MEMBERSHIP  STATUS  AVAILABILITY  MANAGER STATUS
02vfnll4gyc2pg6wvaoweg8hx    worker2    Accepted    Ready   Active
bi7feupsorj32zdkj8fbj1ncd    worker1    Accepted    Ready   Active
cmlouhnpzmmlb59526loze8v3 *  manager  Accepted    Ready   Active        Leader
```

```sh
docker@manager:~$ docker node inspect self
[
    {
        "ID": "cmlouhnpzmmlb59526loze8v3",
        "Version": {
            "Index": 10
        },
        "CreatedAt": "2016-06-29T10:34:10.495840365Z",
        "UpdatedAt": "2016-06-29T10:34:10.692442097Z",
        "Spec": {
            "Role": "manager",
            "Membership": "accepted",
            "Availability": "active"
        },
        "Description": {
            "Hostname": "manager",
            "Platform": {
                "Architecture": "x86_64",
                "OS": "linux"
            },
            "Resources": {
                "NanoCPUs": 1000000000,
                "MemoryBytes": 1044250624
            },
            "Engine": {
                "EngineVersion": "1.12.0-rc2",
                "Labels": {
                    "provider": "virtualbox"
                },
                "Plugins": [
                    {
                        "Type": "Volume",
                        "Name": "local"
                    },
                    {
                        "Type": "Network",
                        "Name": "overlay"
                    },
                    {
                        "Type": "Network",
                        "Name": "host"
                    },
                    {
                        "Type": "Network",
                        "Name": "null"
                    },
                    {
                        "Type": "Network",
                        "Name": "bridge"
                    },
                    {
                        "Type": "Network",
                        "Name": "overlay"
                    }
                ]
            }
        },
        "Status": {
            "State": "ready"
        },
        "ManagerStatus": {
            "Leader": true,
            "Reachability": "reachable",
            "Addr": "192.168.99.100:2377"
        }
    }
]
```

```sh
docker@manager:~$ docker swarm inspect
[
    {
        "ID": "1kjr6jy180zjn39gumjojksln",
        "Version": {
            "Index": 11
        },
        "CreatedAt": "2016-06-29T10:34:10.495823346Z",
        "UpdatedAt": "2016-06-29T10:34:10.768584187Z",
        "Spec": {
            "Name": "default",
            "AcceptancePolicy": {
                "Policies": [
                    {
                        "Role": "worker",
                        "Autoaccept": true
                    },
                    {
                        "Role": "manager",
                        "Autoaccept": false
                    }
                ]
            },
            "Orchestration": {
                "TaskHistoryRetentionLimit": 10
            },
            "Raft": {
                "SnapshotInterval": 10000,
                "LogEntriesForSlowFollowers": 500,
                "HeartbeatTick": 1,
                "ElectionTick": 3
            },
            "Dispatcher": {
                "HeartbeatPeriod": 5000000000
            },
            "CAConfig": {
                "NodeCertExpiry": 7776000000000000
            }
        }
    }
]

```

### 服务集群

搭建一个集群服务, 让我们看看 `docker swarm` 的 服务发现 和 负载均衡. 我们看看搭建的效果.

![img](/images/docker-swarm-look1.png)

我们希望我们有一个集群. 里面运行服务. 其中运行1个 `nginx` 服务和1个 `php-fpm` 服务. `php-fpm` 的规模为 **5**

`nginx` 服务负责在前端代理请求, 将 `php` 类型的请求交给后端的 `php-fpm` 处理. 而 `nginx` 并不需要知道具体要将请求发送给哪个 `php-fpm` 容器. 只要配置好服务名, 剩下的交给 `docker swarm` 的服务发现和负载均衡处理就可以了. 而我们的 `php-fpm` 服务只打印 `uname`.

`/c/Users/Administrator/sixbyte/swarm/fpm/www/index.php`

```php
<?php
echo php_uname();
```

`/c/Users/Administrator/sixbyte/swarm/nginx/conf.d/webapp-php.local.com.conf`

```
server {
    listen       80;
    server_name  webapp-php.local.com;

    location / {
        root   /usr/share/nginx/html;
        index  index.html index.htm;
    }

    error_page   500 502 503 504  /50x.html;
    location = /50x.html {
        root   /usr/share/nginx/html;
    }

    location ~ \.php$ {
        root           /data/www;
        # php-fpm 为服务名, php-fpm 通过监听 9000 端口的方式通讯.
        fastcgi_pass   php-fpm:9000;
        fastcgi_index  index.php;
        fastcgi_param  SCRIPT_FILENAME  $document_root$fastcgi_script_name;
        include        fastcgi_params;
    }

}

```

多主机容器网络. 首先需要理解这个. 要跨主机容器通讯, 需要为容器搭建一个网络. 只要在加入了集群网络的服务才会被服务发现.

```sh
docker@manager:~$ docker network create -d overlay test
```

然后是创建服务, 挂载目录, 连接网络.
```sh
docker@manager:~$ docker service create --replicas 1 --name php-fpm --network test -p 9000 -m type=bind,source=/c/Users/Administrator/sixbyte/swarm/fpm/www,target=/data/www,writable=true php:5-fpm

docker@manager:~$ docker service create --replicas 1 --name nginx --network test -p 80:80 -m type=bind,source=/c/Users/Administrator/sixbyte/swarm/nginx/conf.d,target=/etc/nginx/conf.d,writable=true nginx
```

```
docker@manager:~$ docker service ls
ID            NAME     REPLICAS  IMAGE      COMMAND
00fuklovqhjy  nginx    1/1       nginx
f271805sp66e  php-fpm  1/1       php:5-fpm
```

```
docker@manager:~$ docker service tasks nginx
ID                         NAME     SERVICE  IMAGE  LAST STATE       DESIRED STATE  NODE
dqilzdi8fxsqdbrjycjcxbcon  nginx.1  nginx    nginx  Running 3 Seconds  Running        manager
docker@manager:~$ docker service tasks php-fpm
ID                         NAME       SERVICE  IMAGE      LAST STATE       DESIRED STATE  NODE
epfafd5k8bzuc1yzog169rlrl  php-fpm.1  php-fpm  php:5-fpm  PREPARING 48 Seconds  Running        worker1
```

`PREPARING` 状态一般都是在等待下载镜像..`docker swarm` 会尽量使得每个 `task` 都处于 `DESIRED` 状态.

> 一个任务的生命周期的前半生，大概就是 ASSIGNED -> ACCEPTED -> PREPARING -> STARTING -> RUNNING 。

然后设置 `hosts` , `webapp-php.local.com` 对应 `manager` 的ip地址, 通过浏览器打开:

```
http://webapp-php.local.com/index.php
```

看到如下信息: 

```
Linux 3dbb3dd49bb6 4.4.13-boot2docker #1 SMP Mon Jun 13 23:01:58 UTC 2016 x86_64
```

这是 `php_uname()` 返回的信息, 其中 `3dbb3dd49bb6` 为主机的 `hostname`.

显然, 浏览器->nginx->php-fpm的通讯是通的.

然后我们尝试扩大 `php-fpm` 的规模.

```
docker service scale php-fpm=5
```

在用浏览器访问 webapp-php.local.com/index.php 每次访问得到的 `hostname` 都是不同的. `docker swarm` 为服务提供了负载均衡.


### 滚动更新

```
docker service update --image php:7-fpm php-fpm
```

此时, 等待下载镜像, 服务下所有的tasks都会更换镜像.

其中, 在 `docker service create` 时, `--update-delay` 和 `--update-parallelism` 这2个选项都会影响 `update` 的行为.

### Desired state reconciliation (期望状态调和)

当你删除 `php-fpm` 服务下的某个容器时, `docker swarm` 会重新创建一个新的容器, 保证 scale 的一致. 删除节点后, 节点上需要跑的容器会迁移到其他节点上. 试下就知道.

### 存在的问题

ping php-fpm 和 ping nginx 均得到一个10.0.0.x的ip地址, 也就是swarm的dns服务的解析后得到的ip地址, 但是ping却不可达. 
curl php-fpm 得到的是 节点ip:80 的结果.


docker-machine 里的机器突然全部断电后, 重启顺序 manager, workder.... 但是断电前的容器不会被重启, 不会删除, 而是被忽略, 重新构建新的容器来替代. 这样导致很多僵尸容器. 同时发现服务的规模(scale)的数值超出了设置的值.

docker node promote 试过得到 unreachable 的结果. 但是在manager进程已认为被推举的worker已经是manger, 导致该节点卡死, 并且无法移除. 重启 docker demaon 也没用.只能重新搭建集群...



### 结论

目前 `Swarm` 还没有提供多样的 调度 模型, 也没有完善的健康检查, 监控的特性. 而且目前问题很多, 还不适合在生产环境上使用.