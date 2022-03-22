/*
@Time    : 3/15/22 21:44
@Author  : nil
@File    : tcp.go
*/

package tcp

import (
	"go-proxy/v1/network"
	"go-proxy/v1/socks"
	"log"
	"net"
)

// TcpLocal create a socks server listen on localAddr,
// and this socks server will proxy to remote server.
// localAddr <---> server
func TcpLocal(localAddr, server string, shadow func(net.Conn) net.Conn, socks *socks.Socks) {
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

			lrConn, err := net.Dial("tcp", server)
			if err != nil {
				log.Printf("failed to connect to server %v: %v", server, err)
				return
			}
			defer lrConn.Close()

			lrConn = shadow(lrConn)

			addrByte := network.AddrStrToBytes(tgt)

			if _, err = lrConn.Write(addrByte); err != nil {
				log.Printf("failed to send target address: %v", err)
				return
			}

			log.Printf("proxy %s <-> %s <-> %s", lConn.RemoteAddr(), server, tgt)

			if err = network.Relay(lrConn, lConn); err != nil {
				log.Printf("relay error: %v", err)
			}
		}()
	}
}

// TcpRemote create a relay server listen on addr,
// and this relay server will proxy to target server.
// server <---> target
func TcpRemote(addr string, shadow func(net.Conn) net.Conn) {
	listener, err := net.Listen("tcp", addr)
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

			lrConn := shadow(lrConn)

			tgt, err := network.ReadAddr(lrConn)
			if err != nil {
				log.Printf("failed to get target address from %v: %v", lrConn.RemoteAddr(), err)
				return
			}

			addr := network.AddrBytesToStr(tgt)
			rtConn, err := net.Dial("tcp", addr)

			if err != nil {
				log.Printf("failed to connect to target: %v", err)
				return
			}

			log.Printf("proxy %s <-> %s", lrConn.RemoteAddr(), addr)

			if err = network.Relay(lrConn, rtConn); err != nil {
				log.Printf("relay error: %v", err)
			}
		}()
	}
}


