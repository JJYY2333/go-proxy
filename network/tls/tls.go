/*
@Time    : 3/18/22 22:13
@Author  : Neil
@File    : tls.go
*/

package tls

import (
	sysTLS "crypto/tls"
	"go-proxy/v1/common/statistics"
	"go-proxy/v1/network/util"
	"go-proxy/v1/socks"
	"log"
	"net"
)

func TlsLocal(localAddr, server string, socks *socks.Socks, pair *KeyPair) {
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
			session, err := socks.HandShake(lConn)
			tgt := session.GetTarget()
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
			lrConn, err := sysTLS.Dial("tcp", server, conf)
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

			if _, err = util.Relay(lrConn, lConn); err != nil {
				log.Printf("relay error: %v", err)
			}
		}()
	}
}

func TlsRemote(addr string, clientKeyPair *KeyPair, serverKeyPair *KeyPair) {
	conf, err := GetServerConfig(clientKeyPair, serverKeyPair)
	if err != nil {
		log.Printf("get server tls config error: %v", err)
		return
	}

	listener, err := sysTLS.Listen("tcp", addr, conf)
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

			if _, err := util.Relay(lrConn, rtConn); err != nil {
				log.Printf("relay error: %v", err)
			}
		}()
	}
}

// TlsSolo combine some feature from Local and Remote, so there will be only one proxy server
func TlsSolo(addr string, socks *socks.Socks, clientKeyPair *KeyPair, serverKeyPair *KeyPair) {
	conf, err := GetServerConfig(clientKeyPair, serverKeyPair)
	if err != nil {
		log.Printf("get server tls config error: %v", err)
		return
	}

	listener, err := sysTLS.Listen("tcp", addr, conf)
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

			session, err := socks.HandShake(cConn)
			if err != nil {
				log.Printf("failed to handshake with client %v: %v", cConn.RemoteAddr(), err)
				return
			}
			tgt := session.GetTarget()

			tConn, err := net.Dial("tcp", tgt)

			if err != nil {
				log.Printf("failed to connect to target %s: %v", addr, err)
				return
			}

			log.Printf("proxy between %s <-> %s", addr, tConn.RemoteAddr())

			var n int64
			if n, err = util.Relay(cConn, tConn); err != nil {
				log.Printf("relay error: %v", err)
			}
			statistics.BytesAccount.Add(session.GetUname(), n)
		}()
	}
}
