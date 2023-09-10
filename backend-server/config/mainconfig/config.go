package mainconfig

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
	"github.com/vamika-digital/wms-api-server/config/appconfig"
	"github.com/vamika-digital/wms-api-server/config/authconfig"
	"github.com/vamika-digital/wms-api-server/config/dbconfig"
	"github.com/vamika-digital/wms-api-server/config/serverconfig"
)

var (
	once sync.Once
)

func InitConfig() {
	once.Do(func() {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to get user home directory: %v", err)
		}

		viper.SetConfigName("config.yaml")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.AddConfigPath("config/")    // Path to look for the config file
		viper.AddConfigPath("/etc/wms/")  // Path to look for the config file
		viper.AddConfigPath("$HOME/.wms") // Call multiple times to add many search paths
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err == nil {
			log.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			log.Fatalf("Error reading config file, %s", err)
		}

		// Initialize the sub-configurations
		appconfig.InitAppConfig()
		dbconfig.InitDBConfig()
		serverconfig.InitServerConfig()
		authconfig.InitAuthConfig()
	})
}
