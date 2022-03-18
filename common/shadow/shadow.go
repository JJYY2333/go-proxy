package shadow

import (
	"net"
	"errors"

	"crypto"
	"golang.org/x/crypto/chacha20poly1305"
	// "golang.org/x/crypto/hkdf"
)

// ---------------dummy cipher------------------
type dummy struct {}

func (dummy) StreamConn(con net.Conn) net.Conn{
	return con
}

func (dummy) PacketConn(con net.Conn) net.Conn{
	return con
}

// -------------AEAD cipher-----------------

// AEAD cipher
const (
	aeadAes128Gcm        = "AEAD_AES_128_GCM"
	aeadAes256Gcm        = "AEAD_AES_256_GCM"
	aeadChacha20Poly1305 = "AEAD_CHACHA20_POLY1305"
)



type Cipher interface {
	PacketConnCipher
	StreamConnCipher
}

type PacketConnCipher interface {
	PacketConn(net.PacketConn) net.PacketConn
}

type StreamConnCipher interface {
	StreamConn(net.Conn) net.Conn
}