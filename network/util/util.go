/*
@Time    : 3/20/22 11:34
@Author  : Neil
@File    : connection_util
*/

package util

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

//func ReadAddr(conn net.Conn) ([]byte, error) {
//	buf := make([]byte, 10)
//
//	n, err := io.ReadFull(conn, buf[:6])
//	if n != 6 || err != nil {
//		return nil, fmt.Errorf("read addr error in ReadAddr: %v", err)
//	}
//
//	return buf[:6], nil
//}

// Relay copies between left and right bidirectionally,
// and Return the number of bytes it transfers.
func Relay(left, right net.Conn) (int64, error) {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		var n int64
		n, err1 = io.Copy(right, left)
		log.Printf("=>=>=>: %v bytes flow %v ~~~> %v ~~~> %v", n, left.RemoteAddr(), right.LocalAddr(), right.RemoteAddr())
		right.SetReadDeadline(time.Now().Add(wait)) // unblock read on right
	}()

	var n int64
	n, err = io.Copy(left, right)
	log.Printf("<=<=<=: %v bytes flow %v <~~~ %v <~~~ %v", n, left.RemoteAddr(), left.LocalAddr(), right.RemoteAddr())
	left.SetReadDeadline(time.Now().Add(wait)) // unblock read on left
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return n, err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return n, err
	}
	return n, nil
}

func AddrBytesToStr(b []byte) string {
	ipByte := b[:4]
	portByte := b[4:]
	port := strconv.Itoa(int(binary.BigEndian.Uint16(portByte)))
	ip := net.IP(ipByte).String()
	addr := net.JoinHostPort(ip, port)
	return addr
}

func AddrStrToBytes(s string) []byte {
	ip, port, _ := net.SplitHostPort(s)
	ipByte := []byte(net.ParseIP(ip).To4())
	p, _ := strconv.Atoi(port)
	pByte := make([]byte, 2)
	binary.BigEndian.PutUint16(pByte, uint16(p))
	addrByte := append(ipByte, pByte...)
	return addrByte
}
