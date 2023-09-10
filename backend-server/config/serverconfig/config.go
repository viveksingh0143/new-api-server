package serverconfig

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Address string
	Port    int
}

var Cfg *ServerConfig

func InitServerConfig() {
	Cfg = &ServerConfig{
		Address: viper.GetString("restserver.address"),
		Port:    viper.GetInt("restserver.port"),
	}
}
