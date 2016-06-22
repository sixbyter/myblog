## 12 May 16 一个服务端项目多人开发的 Git 分支模型

-  作者: [sixbyte](http://sixbyte.me/)
-  邮箱: liu.sixbyte@gmail.com
-  禁止转载

公司的项目要迭代一个版本, 我们一帮人一起开发. 分支策略暂时没有一个可用标准, 同时又有新的开发分支要并行开发. 我们的合并各种担惊受怕, 发布到线上的时候不小心回滚了别人的代码2, 3次. 所以, 我们研究了一套可用的 `Git` 分支模型.

要求多个develop分支, 多个hotfixs分支, 也能保证在测试环境, 灰度环境, 生产环境上能有清晰的 network 线路. 


### 废话不多说, 来图

![Git Branch Model](/images/git-branch-model1.png)


### 文档

1. 命名规范
   - 开发环境分支: `develop/*` 比如: `develop/a`, `develop/b`
   - hotfixes分支: `issue/*` 比如: `issue/1`, `issue/2`
   - tag: 年月/日时分. (你肯定会说我2..但是我们是服务端, 用客户端的命名法反而让我们觉得奇怪.) 比如: `201605/201920`
   - 灰度环境分支: `release/*` 比如: `release/develop-a`

2. 说明
   - `online`分支为线上分支
   - `develop/*`分支在测试环境使用, `develop/*`,`issue/*`占用测试环境是竞争关系
   - `release/*`分支从`develop/*`或者`issue/*`上`checkout`, `git checkout -b release/develop-a develop/a`

3. 开发分支或者hotfixes分支上线流程同理. 我以`issue/1`为例说明:
   `issue/1`开发完后, 测试环境`checkout`出`issue/1`分支进行测试. 

   测试环境通过后, 创建新的分支`release/issue-1`. 

   ```shell
   git checkout -b release/issue-1 issue/1
   ```

   灰度环境发布`release/issue-1`最新版本进行测试.问题修复直接提交到`release/issue-1`. 

   测试通过后, `release/issue-1`合并到`online`分支和`issue/1`.

   此时, 请注意! `online`分支必须和`release/issue-1`分支在同一版本号上,  否则, 灰度环境要发布`online`的版本号(假设为`afb`)进行测试. 

   测试通过后, 线上环境发布`afb`, 同时打上`tag`.

   最后`online`广播所以远程分支, 请求合并. 删除`release/issue-1`
