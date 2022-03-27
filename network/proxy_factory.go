/*
@Time    : 3/24/22 21:08
@Author  : Neil
@File    : proxy.go
*/

package network

import (
	"fmt"
	"go-proxy/v1/common/config"
	"go-proxy/v1/network/proxy"
	"go-proxy/v1/network/tcp"
	"go-proxy/v1/network/tls"
	"go-proxy/v1/socks"
	"log"
)

func MakeProxy(cfg *config.Config, socks socks.Socks) (proxy.Proxy, error) {
	var p proxy.Proxy
	var err error

	switch cfg.Connection {
	case "tcp":
		p, err = tcp.NewProxy(cfg, socks)
	case "tls":
		p, err = tls.NewProxy(cfg, socks)
	default:
		err = fmt.Errorf("make proxy error for wrong connection type: %v", cfg.Connection)
	}

	if err != nil {
		return nil, err
	}
	log.Printf("make proxy success, proxy type: %v", cfg.Connection)
	return p, nil
}
