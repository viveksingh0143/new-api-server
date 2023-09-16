package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/vamika-digital/wms-api-server/config/appconfig"
	"github.com/vamika-digital/wms-api-server/config/dbconfig"
	"github.com/vamika-digital/wms-api-server/config/serverconfig"
	"github.com/vamika-digital/wms-api-server/global/drivers"
	"github.com/vamika-digital/wms-api-server/interface/rest"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the web application server",
	Run: func(cmd *cobra.Command, args []string) {
		// Setting timezone to application level
		location, err := time.LoadLocation(appconfig.Cfg.TimeZone)
		if err != nil {
			log.Fatalf("Failed to set time zone: %s", err)
			return
		}
		time.Local = location

		// Create a database connection
		dbConn, err := drivers.NewMySQLConnection(dbconfig.Cfg.Host, dbconfig.Cfg.Port, dbconfig.Cfg.Username, dbconfig.Cfg.Password, dbconfig.Cfg.DBName, dbconfig.Cfg.LogQuery)
		if err != nil {
			log.Fatalf("Failed to connect to database: %s", err)
			return
		}
		// Create a new server instance with the provided port and database connection
		srv := rest.NewServer(serverconfig.Cfg.Address, serverconfig.Cfg.Port, dbConn)
		// Run the server
		srv.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
