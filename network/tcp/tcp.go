/*
@Time    : 3/15/22 21:44
@Author  : nil
@File    : tcp.go
*/

package tcp

import (
	"go-proxy/v1/common/statistics"
	"go-proxy/v1/network/util"
	"go-proxy/v1/socks"
	"log"
	"net"
)

// TcpLocal create a socks server listen on localAddr,
// and this socks server will proxy to remote server.
// localAddr <---> server
func TcpLocal(localAddr, server string, shadow func(net.Conn) net.Conn, socks socks.Socks) {
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
			if err != nil {
				log.Printf("failed to get target address from client: %v", err)
				return
			}
			tgt := session.GetTarget()

			lrConn, err := net.Dial("tcp", server)
			if err != nil {
				log.Printf("failed to connect to server %v: %v", server, err)
				return
			}
			defer lrConn.Close()

			lrConn = shadow(lrConn)

			if _, err = lrConn.Write(tgt); err != nil {
				log.Printf("failed to send target address: %v", err)
				return
			}

			log.Printf("proxy %s <-> %s <-> %s", lConn.RemoteAddr(), server, tgt.String())

			if _, err = util.Relay(lrConn, lConn); err != nil {
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

			tgt, err := util.ReadAddr(lrConn)
			if err != nil {
				log.Printf("failed to get target address from %v: %v", lrConn.RemoteAddr(), err)
				return
			}

			rtConn, err := net.Dial("tcp", tgt.String())

			if err != nil {
				log.Printf("failed to connect to target: %v", err)
				return
			}

			log.Printf("proxy %s <-> %s", lrConn.RemoteAddr(), addr)

			if _, err := util.Relay(lrConn, rtConn); err != nil {
				log.Printf("relay error: %v", err)
			}
		}()
	}
}

func TcpSolo(addr string, shadow func(net.Conn) net.Conn, socks socks.Socks) {
	listener, err := net.Listen("tcp", addr)
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
				log.Printf("failed to get target address from client: %v", err)
				return
			}
			tgt := session.GetTarget()

			tConn, err := net.Dial("tcp", tgt.String())

			if err != nil {
				log.Printf("failed to connect to target %s: %v", addr, err)
				return
			}

			log.Printf("proxy %s <-> %s", tConn.RemoteAddr(), addr)

			var n int64
			if n, err = util.Relay(cConn, tConn); err != nil {
				log.Printf("relay error: %v", err)
			}
			statistics.BytesAccount.Add(session.GetUname(), n)
		}()
	}
}
