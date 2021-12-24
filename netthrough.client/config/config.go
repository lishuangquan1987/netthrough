package config

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

var Config ConfigStruct

func init() {
	cfg, err := ini.Load("netthrough.client.ini")
	if err != nil {
		fmt.Printf("fail to read file:%v", err)
		os.Exit(1)
	}
	Config = ConfigStruct{
		ServerIp:   cfg.Section("client").Key("server_ip").String(),
		ServerPort: cfg.Section("client").Key("server_port").MustInt(10000),
		SourceAddr: cfg.Section("client").Key("source_addr").String(),
	}
}

type ConfigStruct struct {
	ServerIp   string
	ServerPort int
	SourceAddr string
}
