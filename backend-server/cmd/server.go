package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/vamika-digital/wms-api-server/config"
	"github.com/vamika-digital/wms-api-server/global/drivers"
	"github.com/vamika-digital/wms-api-server/interface/rest"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the web application server",
	Run: func(cmd *cobra.Command, args []string) {
		// Setting timezone to application level
		location, err := time.LoadLocation(config.AppConfig.Application.TimeZone)
		if err != nil {
			log.Fatalf("Failed to set time zone: %s", err)
			return
		}
		time.Local = location

		// Create a database connection
		dbConn, err := drivers.NewMySQLConnection(config.AppConfig.Database.Host, config.AppConfig.Database.Port, config.AppConfig.Database.Username, config.AppConfig.Database.Password, config.AppConfig.Database.DBName)
		if err != nil {
			log.Fatalf("Failed to connect to database: %s", err)
			return
		}
		// Create a new server instance with the provided port and database connection
		srv := rest.NewServer(config.AppConfig.RestServer.Address, config.AppConfig.RestServer.Port, dbConn)
		// Run the server
		srv.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
