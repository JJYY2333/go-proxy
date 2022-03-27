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

type Config struct {
	CfgPath    string
	LocalAddr  string
	RemoteAddr string
	ListenAddr string
	UseAuth    bool
	Mode       string
	Connection string
	CertsPath  string
	Shadow     string
	Socks      string
}

func New() *Config {
	cfg := new(Config)
	return cfg
}

func (cfg *Config) LoadConfigFromFile(path string) {
	cfgFile, err := ini.Load(path)
	if err != nil {
		log.Fatalf("fail to parse '%s' : %v", path, err)
	}

	cfg.CfgPath = path

	cfg.loadNet(cfgFile)
	cfg.loadApp(cfgFile)
	cfg.loadLocal(cfgFile)
	cfg.loadRemote(cfgFile)
}

func (cfg *Config) loadNet(cfgFile *ini.File) {
	sec, err := cfgFile.GetSection("net")
	if err != nil {
		log.Fatalf("fail to get section 'server' : %v", err)
	}

	cfg.Connection = sec.Key("Connection").MustString("tcp")
}

func (cfg *Config) loadApp(cfgFile *ini.File) {
	sec, err := cfgFile.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app' : %v", err)
	}

	cfg.UseAuth = sec.Key("Use_Auth").MustBool(false)
	cfg.Mode = sec.Key("Mode").MustString("test")
	cfg.CertsPath = sec.Key("TLS_Certs_Root").MustString("")
	cfg.Shadow = sec.Key("Shadow").MustString("")
	cfg.Socks = sec.Key("Socks").MustString("socksv5")
}

func (cfg *Config) loadLocal(cfgFile *ini.File) {
	sec, err := cfgFile.GetSection("local")
	if err != nil {
		log.Fatalf("Fail to get section 'app' : %v", err)
	}

	cfg.LocalAddr = sec.Key("Local_Address").MustString("127.0.0.1:1089")
	cfg.RemoteAddr = sec.Key("Remote_Address").MustString("127.0.0.1:1090")
}

func (cfg *Config) loadRemote(cfgFile *ini.File) {
	sec, err := cfgFile.GetSection("remote")
	if err != nil {
		log.Fatalf("Fail to get section 'app' : %v", err)
	}

	cfg.ListenAddr = sec.Key("Listen_Address").MustString(":1090")
}
