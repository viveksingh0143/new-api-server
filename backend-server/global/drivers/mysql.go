package drivers

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

type MySQLConnection struct {
	Conn *sqlx.DB
}

func (m *MySQLConnection) GetDB() *sqlx.DB {
	return m.Conn
}

func NewMySQLConnection(host string, port int, username string, password string, dbName string, logQuery bool) (Connection, error) {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		username,
		password,
		host,
		port,
		dbName,
	)

	var conn *sqlx.DB
	var connSqlDB *sql.DB
	var err error

	if logQuery {
		db, err := sql.Open("mysql", connectionString)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, err
		}
		// initiate zerolog
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		zlogger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		// prepare logger
		loggerOptions := []sqldblogger.Option{
			sqldblogger.WithSQLQueryFieldname("sql"),
			sqldblogger.WithWrapResult(false),
			sqldblogger.WithExecerLevel(sqldblogger.LevelDebug),
			sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug),
			sqldblogger.WithPreparerLevel(sqldblogger.LevelDebug),
		}
		// wrap *sql.DB to transparent logger
		connSqlDB = sqldblogger.OpenDriver(connectionString, db.Driver(), zerologadapter.New(zlogger), loggerOptions...)
	}

	if connSqlDB != nil {
		conn = sqlx.NewDb(connSqlDB, "mysql")
	} else {
		// Build the connection string from the config Open a connection to the MySQL database
		conn, err = sqlx.Open("mysql", connectionString)
	}
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	// Test the connection to ensure it's working
	err = conn.Ping()
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return &MySQLConnection{Conn: conn}, nil
}
