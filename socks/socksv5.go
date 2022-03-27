/*
@Time    : 3/26/22 08:49
@Author  : Neil
@File    : socksv51.go
*/

package socks

import (
	"fmt"
	"go-proxy/v1/common/auth"
	"go-proxy/v1/network/util"
	"io"
	"net"
)

type SocksV struct {
	useAuth bool
	checker auth.Authenticator
	session *Session
}

func NewSocksV(use bool, checker auth.Authenticator) *SocksV {
	s := &SocksV{useAuth: use, checker: checker, session: NewSession()}
	return s
}

// HandShake complete the negotiation and returns the target address.
func (s *SocksV) HandShake(conn net.Conn) (*Session, error) {
	err := s.auth(conn)
	if err != nil {
		return nil, fmt.Errorf("socks handshake error: %v", err)
	}

	err = s.connect(conn)
	if err != nil {
		return nil, fmt.Errorf("socks connect error: %v", err)
	}

	return s.session, nil
}

func (s *SocksV) auth(conn net.Conn) error {

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

	//reply ok to client
	//no auth
	if !s.useAuth {
		u := NewAnonymousUser()
		s.session.AddUser(u)
		n, err = conn.Write([]byte{0x05, 0x00})
		if n != 2 || err != nil {
			return fmt.Errorf("write error in Socks5Auth: %v", err)
		}
		return nil
	}

	// reply: use auth
	n, err = conn.Write([]byte{0x05, 0x02})
	if n != 2 || err != nil {
		return fmt.Errorf("write error in Socks5Auth: %v", err)
	}
	// parse uname and password
	n, err = io.ReadFull(conn, buf[:2])
	ver, uLen := buf[0], buf[1]

	n, err = io.ReadFull(conn, buf[:uLen])
	uname := string(buf[:uLen])

	n, err = io.ReadFull(conn, buf[:1])
	pLen := buf[0]
	n, err = io.ReadFull(conn, buf[:pLen])
	passwd := string(buf[:pLen])

	if err != nil {
		return fmt.Errorf("read uname and passwd error:%v", err)
	}

	if s.checker.Check(uname, passwd) {
		u := NewAuthUser(uname)
		s.session.AddUser(u)
		conn.Write([]byte{0x01, 0x00})
		return nil
	} else {
		conn.Write([]byte{0x01, 0x01})
		return fmt.Errorf("auth failed for user: %v", uname)
	}
}

func (s *SocksV) connect(conn net.Conn) error {
	buf := make([]byte, 256)

	// read connect header
	n, err := io.ReadFull(conn, buf[:3])
	if n != 3 {
		return fmt.Errorf("read header error: %v", err)
	}

	ver, cmd, _ := buf[0], buf[1], buf[2]
	if ver != 5 || cmd != 1 {
		return fmt.Errorf("invalid ver/cmd: %d, %d", ver, cmd)
	}

	addr, err := util.ReadAddr(conn)
	if err != nil {
		return err
	}
	s.session.AddTarget(addr)
	n, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	return nil
}
