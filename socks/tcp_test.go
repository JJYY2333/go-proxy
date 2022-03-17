package socks

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestConnection(t *testing.T) {

	dummy := func(conn net.Conn) net.Conn {
		return conn
	}

	go tcpLocal(":1089", ":1090", dummy)

	go tcpRemote(":1090", dummy)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
