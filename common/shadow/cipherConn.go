package shadow

import (
	"net"
	"errors"
	"strings"
	"encoding/hex"
)

type ConnCipher interface {
	PacketConnCipher
	StreamConnCipher
}

type PacketConnCipher interface {
	PacketConn(net.PacketConn) net.PacketConn
}

type StreamConnCipher interface {
	StreamConn(net.Conn) net.Conn
}

// ---------------dummy cipher------------------
type dummy struct{}

func (dummy) StreamConn(c net.Conn) net.Conn             { return c }
func (dummy) PacketConn(c net.PacketConn) net.PacketConn { return c }

//-------------key
func GetKeyByDecode(stringKey string) []byte{
	key, _:= hex.DecodeString(stringKey)
	return key
}

// hardcode的key
var Key = map[int]struct{
	keySize int
	keyString string
	getKey func(string) []byte
}{
	16: {16, StringKeyOf16, GetKeyByDecode},
	32: {32, StringKeyOf32, GetKeyByDecode},
}
// 用的时候，a.getKey(a.keyString)

// 设置全局随机数
var globalNonce = make([] byte, 16)

// CipherConn
const(
	StringKeyOf16 = "6368616e676520746869732070617373"
	StringKeyOf32 = "6368616e676520746869732070617373776f726420746f206120736563726574"

	// AEAD cipher, 可选加密算法
	aeadAes128Gcm        = "AEAD_AES_128_GCM"
	aeadAes256Gcm        = "AEAD_AES_256_GCM"
	aeadChacha20Poly1305 = "AEAD_CHACHA20_POLY1305"

)

// ErrCipherNotSupported occurs when a cipher is not supported (likely because of security concerns).
var ErrCipherNotSupported = errors.New("cipher not supported")

// PickCipher returns a ConnCipher of the given name. Derive key from password if given key is empty.
func PickConnCipher(name string) (ConnCipher, error) {
	name = strings.ToUpper(name)
	length := 0
	switch name {
	case "DUMMY":
		return &dummy{}, nil
	case "CHACHA20-IETF-POLY1305":
		name = aeadChacha20Poly1305
		length = 32
	case "AES-128-GCM":
		name = aeadAes128Gcm
		length = 16
	case "AES-256-GCM":
		name = aeadAes256Gcm
		length = 32
	}

	// key = []byte
	if choice, ok := aeadList[name]; ok {
		//if len(key) != choice.KeySize {
		//	return nil, shadowaead.KeySizeError(choice.KeySize)
		//}
		keyType := Key[length]
		key := keyType.getKey(keyType.keyString)
		aead, err := choice.New(key)
		return &aeadCipher{aead}, err
	}

	return nil, ErrCipherNotSupported
}

type aeadCipher struct{ Cipher }

/*
// 返回加密链接
func (aead *aeadCipher) StreamConn(c net.Conn) net.Conn { return NewStreamConn(c, aead) }
func (aead *aeadCipher) PacketConn(c net.PacketConn) net.PacketConn {
	return NewPacketConn(c, aead)
}
 */