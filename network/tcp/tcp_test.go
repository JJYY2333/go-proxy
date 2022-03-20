package tcp

import (
	"go-proxy/v1/common/auth"
	"go-proxy/v1/socks"
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
	socks := socks.NewSocks(true, auth.NewDummyAuth())
	go TcpLocal(":1089", ":1090", dummy, socks)

	go TcpRemote(":1090", dummy)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
