/*
@Time    : 3/18/22 22:57
@Author  : Neil
@File    : tls_test.go
*/

package tls

import (
	"go-proxy/v1/common/auth"
	"go-proxy/v1/socks"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestConnection(t *testing.T) {

	socks := socks.NewSocks(true, auth.NewDummyAuth())
	go TLSLocal(":1089", ":1090", socks)

	go TLSRemote(":1090")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

func TestServer(t *testing.T) {
	go TLSRemote(":1090")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

func TestLocal(t *testing.T) {
	//go TLSLocal(":1089", "45.76.195.197:443")
	socks := socks.NewSocks(true, auth.NewDummyAuth())
	go TLSLocal(":1089", ":1090", socks)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

func TestTLSSolo(t *testing.T) {
	socks := socks.NewSocks(true, auth.NewDummyAuth())
	go TLSSolo(":1090", socks)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
