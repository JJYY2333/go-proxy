/*
@Time    : 3/27/22 10:13
@Author  : Neil
@File    : socks.go
*/

package socks

import "net"

type Socks interface {
	HandShake(conn net.Conn) (*Session, error)
}
