/*
@Time    : 3/20/22 11:34
@Author  : Neil
@File    : connection_util
*/

package network

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func ReadAddr(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 10)

	n, err := io.ReadFull(conn, buf[:6])

	if n != 6 || err != nil {
		return nil, fmt.Errorf("read addr error in ReadAddr: %v", err)
	}

	return buf[:6], nil
}

// relay copies between left and right bidirectionally
func Relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		var n int64
		n, err1 = io.Copy(right, left)
		log.Printf("%v bytes from %v -> %v", n, left.LocalAddr(), right.LocalAddr())
		right.SetReadDeadline(time.Now().Add(wait)) // unblock read on right
	}()

	var n int64
	n, err = io.Copy(left, right)
	log.Printf("%v bytes from %v -> %v", n, right.LocalAddr(), left.LocalAddr())
	left.SetReadDeadline(time.Now().Add(wait)) // unblock read on left
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}
	return nil
}