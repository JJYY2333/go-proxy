/*
@Time    : 3/17/22 22:47
@Author  : Neil
@File    : main.go
*/

package main

import (
	"flag"
	"go-proxy/v1/common/auth"
	"go-proxy/v1/common/config"
	"go-proxy/v1/network/tcp"
	"go-proxy/v1/network/tls"
	"go-proxy/v1/socks"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var path = flag.String("path", "conf/app_test.ini", "set the ini file path")

func main() {

	flag.Parse()

	dummy := func(conn net.Conn) net.Conn {
		return conn
	}

	cfg := config.New()
	cfg.LoadConfigFromFile(*path)

	socks := socks.NewSocks(cfg.UseAuth, auth.NewDummyAuth())

	if cfg.Connection == "tcp" {
		if cfg.Mode == "local" {
			go tcp.TcpLocal(cfg.LocalAddr, cfg.RemoteAddr, dummy, socks)
		}else if cfg.Mode == "remote" {
			go tcp.TcpRemote(cfg.ListenAddr, dummy)
		}else if cfg.Mode == "test" {
			go tcp.TcpLocal(cfg.LocalAddr, cfg.RemoteAddr, dummy, socks)
			go tcp.TcpRemote(cfg.ListenAddr, dummy)
		} else {
			log.Fatalf("invalid mode in config: %v", cfg.Mode)
		}
	} else if cfg.Connection == "tls" {
		if cfg.Mode == "local" {
			go tls.TLSLocal(cfg.LocalAddr, cfg.RemoteAddr, socks)
		}else if cfg.Mode == "remote" {
			go tls.TLSRemote(cfg.ListenAddr)
		}else if cfg.Mode == "test" {
			go tls.TLSLocal(cfg.LocalAddr, cfg.RemoteAddr, socks)
			go tls.TLSRemote(cfg.ListenAddr)
		} else {
			log.Fatalf("invalid mode in config: %v", cfg.Mode)
		}
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
