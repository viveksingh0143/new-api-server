package dbconfig

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

var Cfg *DBConfig

func InitDBConfig() {
	Cfg = &DBConfig{
		Driver:   viper.GetString("database.driver"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.dbname"),
	}
}

func (c *DBConfig) GetDatabaseConnectionString() string {
	switch c.Driver {
	case "mysql":
		// mysql://user:password@tcp(host:port)/dbname?query
		return fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Username,
			c.Password,
			c.Host,
			c.Port,
			c.DBName,
		)
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?query&sslmode=disable",
			c.Username,
			c.Password,
			c.Host,
			c.Port,
			c.DBName,
		)
	case "sqlite3":
		return c.DBName // SQLite3 might only need the path to the database file
	default:
		log.Fatalf("Unsupported database driver: %s", c.Driver)
		return ""
	}
}
