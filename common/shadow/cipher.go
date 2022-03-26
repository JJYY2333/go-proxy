package shadow

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	//"net"

	//"errors"
	"strconv"
	//"crypto"
	"golang.org/x/crypto/chacha20poly1305"
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

//const(
//	StringKeyOf16 = "6368616e676520746869732070617373"
//	StringKeyOf32 = "6368616e676520746869732070617373776f726420746f206120736563726574"
//
//	// AEAD cipher, 可选加密算法
//	aeadAes128Gcm        = "AEAD_AES_128_GCM"
//	aeadAes256Gcm        = "AEAD_AES_256_GCM"
//	aeadChacha20Poly1305 = "AEAD_CHACHA20_POLY1305"
//)

// -------------AEAD cipher-----------------

/*
生成密钥
keySize := 16, 24, 32
randKey := make([]byte, keySize)
io.ReadFull(rand.Reader, randKey)
fmt.Println(base64.URLEncoding.EncodeToString(key))
key, _ := hex.DecodeString(randKey)
*/

/*
KeySizeError的写法:

type KeySizeError int

func (e KeySizeError) Error() string {
	return "key size error: need " + strconv.Itoa(int(e)) + " bytes"
}
*/

// List of AEAD ciphers: key size in bytes and constructor
var aeadList = map[string]struct {
	KeySize int
	New     func([]byte) (Cipher, error)
}{
	aeadAes128Gcm:        {16, AESGCM},
	aeadAes256Gcm:        {32, AESGCM},
	aeadChacha20Poly1305: {32, Chacha20Poly1305},
}

// 定义一个加密算子应该有的基本功能
type Cipher interface {
	KeySize() int
	// 第一期中不会使用salt对key进行二次hash，留一个接口
	SaltSize() int
	// 如果没有用salt升级key就用key，否则用salt
	Encrypter(key []byte) (cipher.AEAD, error)
	Decrypter(key []byte) (cipher.AEAD, error)
}

// KeySizeError
type KeySizeError int

func (e KeySizeError) Error() string {
	return "key size error: need " + strconv.Itoa(int(e)) + " bytes"
}

type metaCipher struct {
	// psk就是key
	psk      []byte
	makeAEAD func(key []byte) (cipher.AEAD, error)
}

func (a *metaCipher) KeySize() int { return len(a.psk) }
func (a *metaCipher) SaltSize() int {
	if ks := a.KeySize(); ks > 16 {
		return ks
	}
	return 16
}
func (a *metaCipher) Encrypter(salt []byte) (cipher.AEAD, error) {
	//salt=key, 下面的操作都是根据salt和subkey对key进行二次hash，生成新的key
	//subkey := make([]byte, a.KeySize())
	//hkdfSHA1(a.psk, salt, []byte("ss-subkey"), subkey)
	//
	finalKey := processKey(salt, a.psk)
	return a.makeAEAD(finalKey)
}
func (a *metaCipher) Decrypter(salt []byte) (cipher.AEAD, error) {
	//subkey := make([]byte, a.KeySize())
	//hkdfSHA1(a.psk, salt, []byte("ss-subkey"), subkey)
	finalKey := processKey(salt, a.psk)
	return a.makeAEAD(finalKey)
}

// 扩展接口，用salk对key进行加工
func processKey(salt, key []byte) []byte {
	// 目前不对key做二次加工，salt只是用来摆设
	salt[0]++
	return key
}

func aesGCM(key []byte) (cipher.AEAD, error) {
	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(blk)
}

// AESGCM creates a new Cipher with a pre-shared key. len(psk) must be
// one of 16, 24, or 32 to select AES-128/196/256-GCM.
func AESGCM(key []byte) (Cipher, error) {
	switch l := len(key); l {
	case 16, 24, 32: // AES 128/196/256
	default:
		return nil, aes.KeySizeError(l)
	}

	// 生成新的blk，
	//blk, err := aes.NewCipher(key)
	//if err != nil {
	//	return nil, err
	//}

	return &metaCipher{psk: key, makeAEAD: aesGCM}, nil
}

// Chacha20Poly1305 creates a new Cipher with a pre-shared key. len(psk)
// must be 32.
func Chacha20Poly1305(psk []byte) (Cipher, error) {
	if len(psk) != chacha20poly1305.KeySize {
		return nil, KeySizeError(chacha20poly1305.KeySize)
	}
	return &metaCipher{psk: psk, makeAEAD: chacha20poly1305.New}, nil
}
