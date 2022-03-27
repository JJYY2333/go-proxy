/*
@Time    : 3/25/22 22:20
@Author  : Neil
@File    : shadow.go
*/

package tmpShadow

import (
	"fmt"
	"log"
	"net"
)

var (
	shadowFuncMap map[string]func(net.Conn) net.Conn
)

func init() {
	shadowFuncMap = make(map[string]func(net.Conn) net.Conn)
	RegisterShadow("dummy", myDummy)
	//shadowFuncMap["dummy"] = myDummy
}

func RegisterShadow(name string, f func(net.Conn) net.Conn) {
	if _, ok := shadowFuncMap[name]; ok {
		log.Printf("shadow func: %v already registered, drop this register", name)
		return
	}

	shadowFuncMap[name] = f
}

func GetShadow(name string) (func(net.Conn) net.Conn, error) {
	f, ok := shadowFuncMap[name]
	if !ok {
		return nil, fmt.Errorf("failed to find shadow type: %s", name)
	}

	return f, nil
}

func myDummy(conn net.Conn) net.Conn {
	return conn
}
