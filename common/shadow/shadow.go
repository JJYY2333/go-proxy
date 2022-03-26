package shadow

import (
	"fmt"
	"net"
)

var (
	shadowFuncMap map[string]func(net.Conn) net.Conn
)

func init() {
	shadowFuncMap = make(map[string]func(net.Conn) net.Conn)
	shadowFuncMap["dummy"] = myDummy
}

/* 目前有这几种加密算法
dummy是能用的，其它几种还在测试
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
*/

func GetShadow(name string) (func(net.Conn) net.Conn, error) {
	//f, ok := shadowFuncMap[name]
	//if !ok {
	//	return nil, fmt.Errorf("failed to find shadow type: %s", name)
	//}
	//
	//return f, nil

	f, ok := PickConnCipher(name)
	if ok != nil {
		return nil, fmt.Errorf("failed to find shadow type: %s", name)
	}

	return f.StreamConn, nil
}

func myDummy(conn net.Conn) net.Conn {
	return conn
}
