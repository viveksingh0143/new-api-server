package appconfig

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	TimeZone string
}

var Cfg *AppConfig

func InitAppConfig() {
	Cfg = &AppConfig{
		TimeZone: viper.GetString("application.timezone"),
	}
}
