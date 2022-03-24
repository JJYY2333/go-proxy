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
	"go-proxy/v1/network/tls"
	"go-proxy/v1/socks"
)

func MakeProxy(cfg *config.Config, socks *socks.Socks) (proxy.Proxy, error) {
	switch cfg.Connection {
	case "tcp":
		return nil, fmt.Errorf("currently unsupported")
	case "tls":
		p, err := tls.NewProxy(cfg, socks)
		if err != nil {
			return nil, err
		}
		return p, nil
	default:
		return nil, fmt.Errorf("make proxy error for wrong connection type: %v", cfg.Connection)
	}
}
