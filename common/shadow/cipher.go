package shadow

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net"
	//"errors"
	//
	//"crypto"
	//"golang.org/x/crypto/chacha20poly1305"
	// "golang.org/x/crypto/hkdf"
)

/*
	该模块实现加密算子cipher，根据传入的参数返回对应的cipher算子
	默认由外部调用提供keySize和key参数
*/

// 生成密钥
func KEYGen(keySize int) ([]byte, string) {
	randKey := make([]byte, keySize)
	io.ReadFull(rand.Reader, randKey)
	key := base64.URLEncoding.EncodeToString(randKey)
	// 返回生成的key数组和字符串key
	return randKey, key
}

// ---------------dummy cipher------------------
type dummy struct{}

func (dummy) StreamConn(con net.Conn) net.Conn {
	return con
}

func (dummy) PacketConn(con net.Conn) net.Conn {
	return con
}

// -------------AEAD cipher-----------------

/*
生成密钥
keySize := 16, 24, 32
randKey := make([]byte, keySize)
io.ReadFull(rand.Reader, randKey)
fmt.Println(base64.URLEncoding.EncodeToString(key))
key, _ := hex.DecodeString(randKey)
*/

// AEAD cipher, 可选加密算法
const (
	aeadAes128Gcm        = "AEAD_AES_128_GCM"
	aeadAes256Gcm        = "AEAD_AES_256_GCM"
	aeadChacha20Poly1305 = "AEAD_CHACHA20_POLY1305"
)

// 定义一个加密算子应该有的基本功能
type Cipher interface {
	KeySize() int
	SaltSize() int
	Encrypter(salt []byte) (cipher.AEAD, error)
	Decrypter(salt []byte) (cipher.AEAD, error)
}

/*
KeySizeError的写法:

type KeySizeError int

func (e KeySizeError) Error() string {
	return "key size error: need " + strconv.Itoa(int(e)) + " bytes"
}
*/

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
