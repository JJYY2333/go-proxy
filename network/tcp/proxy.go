/*
@Time    : 3/25/22 21:34
@Author  : Neil
@File    : proxy.go
*/

package tcp

import (
	"fmt"
	"go-proxy/v1/common/config"
	"go-proxy/v1/common/shadow"
	"go-proxy/v1/socks"
	"net"
)

type TcpProxy struct {
	// shadow
	shadow func(net.Conn) net.Conn
	//local
	laddr string
	raddr string

	//remote or solo
	listenAddr string

	//socks
	socks socks.Socks

	//mode
	mode string
}

func NewProxy(cfg *config.Config, socks socks.Socks) (*TcpProxy, error) {
	p := new(TcpProxy)
	p.socks = socks
	p.mode = cfg.Mode
	var err error
	p.shadow, err = shadow.GetShadow(cfg.Shadow)
	if err != nil {
		return nil, fmt.Errorf("build proxy error when getting shadow: %v", err)
	}
	switch p.mode {
	case "local":
		p.laddr = cfg.LocalAddr
		p.raddr = cfg.RemoteAddr
	case "remote":
		p.listenAddr = cfg.ListenAddr
	case "test":
		p.laddr = cfg.LocalAddr
		p.raddr = cfg.RemoteAddr
		p.listenAddr = cfg.ListenAddr
	case "solo":
		p.listenAddr = cfg.ListenAddr
	default:
		return nil, fmt.Errorf("build proxy error for there is no mode in tcp proxy: %v", p.mode)
	}

	return p, nil
}

func (p *TcpProxy) Start() {
	switch p.mode {
	case "local":
		p.startLocal()
	case "remote":
		p.startRemote()
	case "test":
		p.startTest()
	case "solo":
		p.startSolo()
	}
}

func (p *TcpProxy) startLocal() {
	go TcpLocal(p.laddr, p.raddr, p.shadow, p.socks)
}

func (p *TcpProxy) startRemote() {
	go TcpRemote(p.listenAddr, p.shadow)
}

func (p *TcpProxy) startTest() {
	p.startLocal()
	p.startRemote()
}

func (p *TcpProxy) startSolo() {
	go TcpSolo(p.listenAddr, p.shadow, p.socks)
}
