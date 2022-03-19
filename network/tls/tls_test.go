/*
@Time    : 3/18/22 22:57
@Author  : Neil
@File    : tls_test.go
*/

package tls

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestConnection(t *testing.T) {

	go TLSLocal(":1089", ":1090")

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

	go TLSLocal(":1089", ":1090")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}