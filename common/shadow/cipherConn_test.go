package shadow

import (
	"go-proxy/v1/common/auth"
	"go-proxy/v1/network/tcp"
	"go-proxy/v1/socks"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestCipherConn(t *testing.T) {
	cipherType := "AEAD_AES_256_GCM"

	cipherConn, err := PickConnCipher(cipherType)

	if err != nil {
		print(err)
	}
	//print(cipherConn)
	socks := socks.NewSocks(true, auth.NewDummyAuth())
	go tcp.TcpLocal(":1089", ":1090", cipherConn.StreamConn, socks)

	go tcp.TcpRemote(":1090", cipherConn.StreamConn)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	/*
		dummy := func(conn net.Conn) net.Conn {
			return conn
		}
		socks := socks.NewSocks(true, auth.NewDummyAuth())
		go TcpLocal(":1089", ":1090", dummy, socks)

		go TcpRemote(":1090", dummy)

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
	*/

}
