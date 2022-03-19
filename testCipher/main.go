//package main
//
//import (
//	cryptorand "crypto/rand"
//	"fmt"
//	"golang.org/x/crypto/chacha20poly1305"
//	"log"
//)
//
//func Cha(key []byte) {
//	aead, err := chacha20poly1305.NewX(key)
//	if err != nil {
//		log.Fatalln("Failed to instantiate XChaCha20-Poly1305:", err)
//	}
//
//	for _, msg := range []string{
//		"Attack at dawn.",
//		"The eagle has landed.",
//		"Gophers, gophers, gophers everywhere!",
//	} {
//		// Encryption.
//		nonce := make([]byte, chacha20poly1305.NonceSizeX)
//		if _, err := cryptorand.Read(nonce); err != nil {
//			panic(err)
//		}
//		ciphertext := aead.Seal(nil, nonce, []byte(msg), nil)
//
//		// Decryption.
//		plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
//		if err != nil {
//			log.Fatalln("Failed to decrypt or authenticate message:", err)
//		}
//
//		fmt.Printf("%s\n", plaintext)
//	}
//}
//
//func main() {
//	var key = make([]byte, chacha20poly1305.KeySize)
//	fmt.Printf("", key)
//	Cha(key)
//}

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

func main() {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	ciphertext, _ := hex.DecodeString("c3aaa29f002ca75870806e44086700f62ce4d43e902b3888e23ceff797a7a471")
	nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")

	fmt.Printf("%x\n", nonce)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	println("block: ", block.BlockSize())

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	println("aesgcm1: ", aesgcm.NonceSize())
	println("aesgcm2: ", aesgcm.Overhead())

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%s\n", plaintext)

}
