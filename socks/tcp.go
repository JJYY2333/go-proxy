/*
@Time    : 3/15/22 21:44
@Author  : nil
@File    : tcp.go
*/

package socks

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func tcpLocal(localAddr, server string, shadow func(net.Conn) net.Conn) {
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Printf("failed to listen on %s: %v\n", localAddr, err)
	}

	for {
		lConn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept: %v\n", err)
			continue
		}

		go func(){
			defer lConn.Close()
			tgt, err := socks5GetAddr(lConn)
			if err != nil {
				log.Printf("failed to get target address from client: %v", err)
			}

			lrConn, err := net.Dial("tcp", server)
			if err != nil {
				log.Printf("failed to connect to server %v: %v", server, err)
			}
			defer lrConn.Close()

			lrConn = shadow(lrConn)

			if _, err = lrConn.Write([]byte(tgt)); err != nil {
				log.Printf("failed to send target address: %v", err)
				return
			}

			log.Printf("proxy %s <-> %s <-> %s", lConn.RemoteAddr(), server, tgt)

			if err = relay(lrConn, lConn); err != nil {
				log.Printf("relay error: %v", err)
			}
		}()
	}
}


func tcpRemote(addr string, shadow func(net.Conn) net.Conn) {
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
		go func(){
			defer lrConn.Close()


			lrConn := shadow(lrConn)

			tgt, err := ReadAddr(lrConn)
			if err != nil {
				log.Printf("failed to get target address from %v: %v", lrConn.RemoteAddr(), err)
				return
			}

			rtConn, err := net.Dial("tcp", tgt.String())

			if err != nil {
				log.Printf("failed to connect to target: %v", err)
			}

			log.Printf("proxy %s <-> %s", lrConn.RemoteAddr(), tgt)

			if err = relay(lrConn, rtConn); err != nil {
				log.Printf("relay error: %v", err)
			}
		}()
	}
}

func ReadAddr(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 10)
	
	n, err := io.Readfull(conn, buf[:7])
	
	if n != 7 || err != nil {
		return nil, fmt.Errorf("read addr error in ReadAddr: %v", err)	
	}
	
	return buf[:7], nil
}

// relay copies between left and right bidirectionally
func relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, left)
		right.SetReadDeadline(time.Now().Add(wait)) // unblock read on right
	}()
	_, err = io.Copy(left, right)
	left.SetReadDeadline(time.Now().Add(wait)) // unblock read on left
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}
	return nil
}


//整个协商， 子协商， 请求阶段， 已拿到target address为结束
func socks5GetAddr(conn net.Conn) (string, error){
	buf := make([]byte, 256)
	n, err := io.ReadFull(conn, buf[:2])
	if n != 2 {
		return "", fmt.Errorf("reading Header error in Socks5Auth: %v", err)
	}

	ver, nMethods := buf[0], buf[1]
	if ver != 5 {
		return "", fmt.Errorf("invalid version: %v", ver)
	}

	n, err = io.ReadFull(conn, buf[:nMethods])
	if n != int(nMethods) {
		return "", fmt.Errorf("reading methods: %v", err)
	}

	n, err = conn.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		return "", fmt.Errorf("write error in Socks5Auth: %v", err)
	}


	n, err = io.ReadFull(conn, buf[:4])
	if n != 4 {
		return "", errors.New("read header: " + err.Error())
	}

	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 5 || cmd != 1 {
		return "", errors.New("invalid ver/cmd")
	}

	addr := ""
	switch atyp {
	case 1:
		n, err = io.ReadFull(conn, buf[:4])
		if n != 4 {
			return "", errors.New("invalid IPv4: " + err.Error())
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])

	case 3:
		n, err = io.ReadFull(conn, buf[:1])
		if n != 1 {
			return "", errors.New("invalid hostname: " + err.Error())
		}
		addrLen := int(buf[0])

		n, err = io.ReadFull(conn, buf[:addrLen])
		if n != addrLen {
			return "", errors.New("invalid hostname: " + err.Error())
		}
		addr = string(buf[:addrLen])

	case 4:
		return "", errors.New("IPv6: no supported yet")

	default:
		return "", errors.New("invalid atyp")
	}

	n, err = io.ReadFull(conn, buf[:2])
	if n != 2 {
		return "", errors.New("read port: " + err.Error())
	}
	port := binary.BigEndian.Uint16(buf[:2])

	destAddrPort := fmt.Sprintf("%s:%d", addr, port)

	return destAddrPort, nil
}
