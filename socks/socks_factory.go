/*
@Time    : 3/27/22 10:11
@Author  : Neil
@File    : socks_factory.go
*/

package socks

import (
	"fmt"
	"go-proxy/v1/common/auth"
	"go-proxy/v1/common/config"
	"log"
)

func MakeSocks(cfg *config.Config) (Socks, error) {
	var s Socks
	var err error
	switch cfg.Socks {
	case "socksv5":
		use := cfg.UseAuth
		checker := auth.NewDummyAuth()
		s = NewSocksV(use, checker)
	case "shadowsocks":
		s = NewShadowSocks()
	default:
		err = fmt.Errorf("make socks error for wrong socks type: %v", cfg.Socks)
	}

	if err != nil {
		return nil, err
	}
	log.Printf("make socks success, socks type: %v", cfg.Socks)
	return s, nil
}
