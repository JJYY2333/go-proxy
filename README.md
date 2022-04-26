# go-proxy
A proxy program based on Go.

support two modes: 
* local server + remote server
* solo

In local + remote mode, local server listenning on a local port on your machine communicates with clients that support socks proxy(eg. safari, chrome web browser). And remote server shoulb be deployed on a machine outside GFW. you can start up these two server by modifying .ini file.

In solo mode, local server and remote server has been merged together, so it can communicates with specific applications on your mobile phone(like shadowsocks rocket). You can also start up a solo server by modifying .ini file.

support two transport layer protocl:
* tcp
* tls

when using tcp, you should define an encrypted algorithm to bypass the GFW. When using tls, you can depends on tls's encrypted algorithm.


## 结构图

![Structure](structure.png)

