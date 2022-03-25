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
