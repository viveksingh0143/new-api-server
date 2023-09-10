package authconfig

import (
	"github.com/spf13/viper"
)

type AuthConfig struct {
	ExpiryDuration     int64
	ExpiryLongDuration int64
	SecretKey          string
}

var Cfg *AuthConfig

func InitAuthConfig() {
	Cfg = &AuthConfig{
		ExpiryDuration:     viper.GetInt64("auth.expiry_duration"),
		ExpiryLongDuration: viper.GetInt64("auth.expiry_long_duration"),
		SecretKey:          viper.GetString("auth.secret_key"),
	}
}
