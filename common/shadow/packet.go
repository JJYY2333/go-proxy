package shadow

import (
	"net"
	"sync"
)

type packetConn struct {
	net.PacketConn
	Cipher
	sync.Mutex
	buf []byte // write lock
}

const maxPacketSize = 64 * 1024

func NewPacketConn(c net.PacketConn, ciph Cipher) net.PacketConn {
	return &packetConn{PacketConn: c, Cipher: ciph, buf: make([]byte, maxPacketSize)}
}
