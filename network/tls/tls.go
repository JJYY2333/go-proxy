/*
@Time    : 3/18/22 22:13
@Author  : Neil
@File    : tls.go
*/

package tls

import (
	"crypto/tls"
	"go-proxy/v1/network/util"
	"go-proxy/v1/socks"
	"log"
	"net"
)

func TLSLocal(localAddr, server string, socks *socks.Socks, pair *KeyPair) {
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Printf("failed to listen on %s: %v", localAddr, err)
	}
	log.Printf("listening TCP on %s", localAddr)

	for {
		lConn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept: %v\n", err)
			continue
		}

		go func() {
			defer lConn.Close()
			tgt, err := socks.HandShake(lConn)
			if err != nil {
				log.Printf("failed to get target address from client: %v", err)
				return
			}

			// set tls config
			conf, err := GetClientConfig(pair)
			if err != nil {
				log.Fatalf("get client tls config error: %v", err)
			}

			// dial tls
			lrConn, err := tls.Dial("tcp", server, conf)
			if err != nil {
				log.Printf("failed to connect to server %v: %v", server, err)
				return
			}
			defer lrConn.Close()

			addrByte := util.AddrStrToBytes(tgt)
			if _, err = lrConn.Write(addrByte); err != nil {
				log.Printf("failed to send target address: %v", err)
				return
			}

			log.Printf("proxy %s <-> %s <-> %s", lConn.RemoteAddr(), server, tgt)

			if err = util.Relay(lrConn, lConn); err != nil {
				log.Printf("relay error: %v", err)
			}
		}()
	}
}

func TLSRemote(addr string, clientKeyPair *KeyPair, serverKeyPair *KeyPair) {
	conf, err := GetServerConfig(clientKeyPair, serverKeyPair)
	if err != nil {
		log.Printf("get server tls config error: %v", err)
		return
	}

	listener, err := tls.Listen("tcp", addr, conf)
	if err != nil {
		log.Printf("failed to listen on %s: %v", addr, err)
		return
	}

	log.Printf("listening TCP on %s", addr)

	for {
		lrConn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept: %v\n", err)
			continue
		}

		//不用担心lrConn全部使用的是一个， 这里用到了闭包
		go func() {
			defer lrConn.Close()

			tgt, err := util.ReadAddr(lrConn)
			if err != nil {
				log.Printf("failed to get target address from %v: %v", lrConn.RemoteAddr(), err)
				return
			}

			addr := util.AddrBytesToStr(tgt)
			//log.Printf("remote tgt is: %v, length is :%v, string is :%v", addr, len(addr), string(addr))

			rtConn, err := net.Dial("tcp", addr)

			if err != nil {
				log.Printf("failed to connect to target %s: %v", addr, err)
				return
			}

			log.Printf("proxy %s <-> %s", lrConn.RemoteAddr(), addr)

			if err = util.Relay(lrConn, rtConn); err != nil {
				log.Printf("relay error: %v", err)
			}
		}()
	}
}

// TLSSolo combine some feature from Local and Remote, so there will be only one proxy server
func TLSSolo(addr string, socks *socks.Socks, clientKeyPair *KeyPair, serverKeyPair *KeyPair) {
	conf, err := GetServerConfig(clientKeyPair, serverKeyPair)
	if err != nil {
		log.Printf("get server tls config error: %v", err)
		return
	}

	listener, err := tls.Listen("tcp", addr, conf)
	if err != nil {
		log.Printf("failed to listen on %s: %v", addr, err)
		return
	}

	log.Printf("listening TCP on %s", addr)

	for {
		cConn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept: %v\n", err)
			continue
		}

		go func() {
			defer cConn.Close()

			tgt, err := socks.HandShake(cConn)
			if err != nil {
				log.Printf("failed to get target address from client: %v", err)
				return
			}

			tConn, err := net.Dial("tcp", tgt)

			if err != nil {
				log.Printf("failed to connect to target %s: %v", addr, err)
				return
			}

			log.Printf("proxy %s <-> %s", tConn.RemoteAddr(), addr)

			if err = util.Relay(cConn, tConn); err != nil {
				log.Printf("relay error: %v", err)
			}
		}()
	}
}
