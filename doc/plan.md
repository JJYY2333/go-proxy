第一步：communication
第二步：shadow

* 实现DNS代理 [***], 三颗星的优先级的理由是，墙内的DNS也并非全部被污染了，还是有很多网站可以用来测试的, 使用8.8.8.8作为DNS服务器也是可行的
* [DONE] 实现TLS作为加密层 [****], 理论上是最容易的加密手段，而且实现起来感觉不麻烦。
* [DONE] [****] 重构socks握手过程
* [****] tcp 与 tls有太多重复的地方， 需要重构， streamConn, tcp or tls? 
* socks over udp [**]
* [***] 可以设置启动配置文件的命令行参数
* [***] 由于各种Panic导致的异常情况的处理， 例如 check ip address, dial timeout 等等。给出session关闭的信息
* [****] 重构tcp, tls，它们应该是一个结构体，这样能使用config初始化
* [***] 完成一个第三方认证的auth


