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
	"go-proxy/v1/network"
	"go-proxy/v1/socks"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var path = flag.String("path", "any/conf/app_test.ini", "set the ini file path")

func main() {

	flag.Parse()

	cfg := config.New()
	cfg.LoadConfigFromFile(*path)

	socks := socks.NewSocks(cfg.UseAuth, auth.NewDummyAuth())
	proxy, err := network.MakeProxy(cfg, socks)
	if err != nil {
		log.Fatalf("make proxy error: %v", err)
	}
	proxy.Start()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
