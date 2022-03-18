/*
@Time    : 3/18/22 22:13
@Author  : Neil
@File    : tls.go
*/

package tls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func TLSLocal(localAddr, server string) {
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
			tgt, err := socks5GetAddr(lConn)
			log.Printf("local tgt :%v", tgt)
			if err != nil {
				log.Printf("failed to get target address from client: %v", err)
			}

			// set tls config
			cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
			if err != nil {
				log.Println(err)
				return
			}
			certBytes, err := ioutil.ReadFile("certs/client.pem")
			if err != nil {
				panic("Unable to read cert.pem")
			}
			clientCertPool := x509.NewCertPool()
			ok := clientCertPool.AppendCertsFromPEM(certBytes)
			if !ok {
				panic("failed to parse root certificate")
			}
			conf := &tls.Config{
				RootCAs:            clientCertPool,
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
			}

			// dial tls
			lrConn, err := tls.Dial("tcp", server, conf)
			if err != nil {
				log.Printf("failed to connect to server %v: %v", server, err)
			}
			defer lrConn.Close()

			ip, port, err := net.SplitHostPort(tgt)
			ip_byte := []byte(net.ParseIP(ip).To4())
			p, err := strconv.Atoi(port)
			p_byte := make([]byte, 2)
			binary.BigEndian.PutUint16(p_byte, uint16(p))
			addr_byte := append(ip_byte, p_byte...)

			if _, err = lrConn.Write(addr_byte); err != nil {
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

func TLSRemote(addr string) {
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Println(err)
		return
	}
	certBytes, err := ioutil.ReadFile("certs/client.pem")
	if err != nil {
		panic("Unable to read cert.pem")
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("failed to parse root certificate")
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}
	listener, err := tls.Listen("tcp", addr, config)
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

			tgt, err := ReadAddr(lrConn)
			if err != nil {
				log.Printf("failed to get target address from %v: %v", lrConn.RemoteAddr(), err)
				return
			}

			ipByte := tgt[:4]
			portByte := tgt[4:]

			port := strconv.Itoa(int(binary.BigEndian.Uint16(portByte)))
			ip := net.IP(ipByte).String()
			addr := net.JoinHostPort(ip, port)
			log.Printf("remote tgt is: %v, length is :%v, string is :%v", addr, len(addr), string(addr))

			rtConn, err := net.Dial("tcp", addr)

			if err != nil {
				log.Printf("failed to connect to target: %v", err)
			}

			log.Printf("proxy %s <-> %s", lrConn.RemoteAddr(), addr)

			if err = relay(lrConn, rtConn); err != nil {
				log.Printf("relay error: %v", err)
			}
		}()
	}
}

func ReadAddr(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 10)

	n, err := io.ReadFull(conn, buf[:6])

	if n != 6 || err != nil {
		return nil, fmt.Errorf("read addr error in ReadAddr: %v", err)
	}

	return buf[:6], nil
}

// relay copies between left and right bidirectionally
func relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		var n int64
		n, err1 = io.Copy(right, left)
		log.Printf("%v bytes from %v -> %v", n, left.LocalAddr(), right.LocalAddr())
		right.SetReadDeadline(time.Now().Add(wait)) // unblock read on right
	}()

	var n int64
	n, err = io.Copy(left, right)
	log.Printf("%v bytes from %v -> %v", n, right.LocalAddr(), left.LocalAddr())
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
func socks5GetAddr(conn net.Conn) (string, error) {
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

	n, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return "", errors.New("write rsp: " + err.Error())
	}
	return destAddrPort, nil
}
