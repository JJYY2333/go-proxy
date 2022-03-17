/*
@Time    : 3/17/22 22:25
@Author  : Neil
@File    : config.go
*/

package config

import (
	"github.com/go-ini/ini"
	"log"
)

var (
	Cfg *ini.File
	LocalAddr string
	RemoteAddr string
	UseAuth bool
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini : %v", err)
	}
	LoadServer()
	LoadApp()
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server' : %v", err)
	}

	LocalAddr = sec.Key("Local_Address").MustString("127.0.0.1:1089")
	RemoteAddr = sec.Key("Remote_Address").MustString("127.0.0.1:1090")
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app' : %v", err)
	}

	UseAuth = sec.Key("Use_Auth").MustBool(false)
}