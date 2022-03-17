/*
@Time    : 3/17/22 22:06
@Author  : Neil
@File    : socksv5
*/

package socks

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

// HandShake complete the negotiation and returns the target address.
func HandShake(conn net.Conn) (string, error) {
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

func auth(conn net.Conn) error {
	buf := make([]byte, 256)
	// read VER, nMETHODS
	n, err := io.ReadFull(conn, buf[:2])
	if n != 2 {
		return fmt.Errorf("reading Header error in Socks5Auth: %v", err)
	}

	//check VERSION
	ver, nMethods := buf[0], buf[1]
	if ver != 5 {
		return fmt.Errorf("invalid version: %v", ver)
	}

	//check METHODS
	n, err = io.ReadFull(conn, buf[:nMethods])
	if n != int(nMethods) {
		return fmt.Errorf("reading methods: %v", err)
	}

}