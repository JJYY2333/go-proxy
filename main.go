/*
@Time    : 3/17/22 22:47
@Author  : Neil
@File    : main.go
*/

package main

import (
	"go-proxy/v1/common/config"
	"go-proxy/v1/network/tcp"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dummy := func(conn net.Conn) net.Conn {
		return conn
	}

	go tcp.TcpLocal(config.LocalAddr, config.RemoteAddr, dummy)
	go tcp.TcpRemote(config.RemoteAddr, dummy)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
