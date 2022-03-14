/*
@Time    : 3/14/22 22:35
@Author  : nil
@File    : main
*/

package main

import (
"encoding/binary"
"errors"
"flag"
"fmt"
"io"
"log"
"net"
)

var auth = flag.Bool("auth", false, "set true if use auth")

func main(){

	flag.Parse()

	listener, err := net.Listen("tcp", ":1089")
	if err != nil {
		log.Printf("listen error in main: %v\n", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept failed in main: %v\n", err)
			continue
		}

		go handler(conn)
	}

}

func handler(conn net.Conn) {
	if err := Socks5Auth(conn); err != nil {
		log.Printf("auth error in handler: %v\n", err)
		conn.Close()
		return
	}

	target, err := Socks5Connect(conn)
	if err != nil {
		log.Printf("connect error in handler: %v\n", err)
		conn.Close()
		return
	}

	Socks5Forward(conn, target)
}


func Socks5Auth(conn net.Conn) error {
	buf := make([]byte, 256)
	n, err := io.ReadFull(conn, buf[:2])
	if n != 2 {
		return fmt.Errorf("reading Header error in Socks5Auth: %v", err)
	}

	ver, nMethods := buf[0], buf[1]
	if ver != 5 {
		return fmt.Errorf("invalid version: %v", ver)
	}

	n, err = io.ReadFull(conn, buf[:nMethods])
	if n != int(nMethods) {
		return fmt.Errorf("reading methods: %v", err)
	}

	fmt.Printf("use auth: %v\n", *auth)
	//no auth
	if !*auth {
		n, err = conn.Write([]byte{0x05, 0x00})
		if n != 2 || err != nil {
			return fmt.Errorf("write error in Socks5Auth: %v", err)
		}
		return nil
	}

	// use auth, parse
	n, err = io.ReadFull(conn, buf[:2])
	ver, ulen := buf[0], buf[1]

	n, err = io.ReadFull(conn, buf[:ulen])
	uname := buf[:ulen]

	n, err = io.ReadFull(conn, buf[:1])
	plen := buf[0]
	n, err = io.ReadFull(conn, buf[:plen])
	passwd := buf[:plen]

	if string(uname) != "neil" && string(passwd) != "hello" {
		return fmt.Errorf("auth failed for user: %v", uname)
	}


	return nil
}

func Socks5Connect(client net.Conn) (net.Conn, error) {
	buf := make([]byte, 256)

	n, err := io.ReadFull(client, buf[:4])
	if n != 4 {
		return nil, errors.New("read header: " + err.Error())
	}

	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 5 || cmd != 1 {
		return nil, errors.New("invalid ver/cmd")
	}

	addr := ""
	switch atyp {
	case 1:
		n, err = io.ReadFull(client, buf[:4])
		if n != 4 {
			return nil, errors.New("invalid IPv4: " + err.Error())
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])

	case 3:
		n, err = io.ReadFull(client, buf[:1])
		if n != 1 {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addrLen := int(buf[0])

		n, err = io.ReadFull(client, buf[:addrLen])
		if n != addrLen {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addr = string(buf[:addrLen])

	case 4:
		return nil, errors.New("IPv6: no supported yet")

	default:
		return nil, errors.New("invalid atyp")
	}

	n, err = io.ReadFull(client, buf[:2])
	if n != 2 {
		return nil, errors.New("read port: " + err.Error())
	}
	port := binary.BigEndian.Uint16(buf[:2])

	destAddrPort := fmt.Sprintf("%s:%d", addr, port)
	dest, err := net.Dial("tcp", destAddrPort)
	if err != nil {
		return nil, errors.New("dial dst: " + err.Error())
	}

	n, err = client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		dest.Close()
		return nil, errors.New("write rsp: " + err.Error())
	}

	return dest, nil
}

func Socks5Forward(client, target net.Conn) {
	forward := func(src, dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}
	go forward(client, target)
	go forward(target, client)
}
