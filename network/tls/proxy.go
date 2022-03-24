/*
@Time    : 3/23/22 22:58
@Author  : Neil
@File    : proxy.go
*/

package tls

import (
	"fmt"
	"go-proxy/v1/common/config"
	"go-proxy/v1/socks"
	"path"
)

type TLSProxy struct {
	clientKeyPair *KeyPair
	serverKeyPair *KeyPair

	//local
	laddr string
	raddr string

	//remote or solo
	listenAddr string

	//socks
	socks *socks.Socks

	//mode
	mode string
}

func NewProxy(cfg *config.Config, socks *socks.Socks) (*TLSProxy, error) {
	p := new(TLSProxy)
	p.clientKeyPair = &KeyPair{key: path.Join(cfg.CertsPath, "client.key"), pem: path.Join(cfg.CertsPath, "client.pem")}
	p.serverKeyPair = &KeyPair{key: path.Join(cfg.CertsPath, "server.key"), pem: path.Join(cfg.CertsPath, "server.pem")}
	p.socks = socks
	p.mode = cfg.Mode

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
		return nil, fmt.Errorf("config Mode error, no mode in tls proxy: %v", p.mode)
	}

	return p, nil
}

// Start build different mode of this proxy
func (p *TLSProxy) Start() {
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

func (p *TLSProxy) startLocal() {
	go TLSLocal(p.laddr, p.raddr, p.socks, p.clientKeyPair)
}

func (p *TLSProxy) startRemote() {
	go TLSRemote(p.listenAddr, p.clientKeyPair, p.serverKeyPair)
}

func (p *TLSProxy) startSolo() {
	go TLSSolo(p.listenAddr, p.socks, p.clientKeyPair, p.serverKeyPair)
}

func (p *TLSProxy) startTest() {
	p.startLocal()
	p.startRemote()
}
